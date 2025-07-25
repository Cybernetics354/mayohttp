package app

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Cybernetics354/mayohttp/app/telescope"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *State) OpenTelescope(msg openTelescopeMsg) (tea.Model, tea.Cmd) {
	var items []list.Item
	var title string

	switch msg.teleType {
	case TELESCOPE_METHOD_PALLETE:
		items = methodPalletesTelescope
		title = "Select Method"
	case TELESCOPE_QUICK_ACCESS:
		items = quickAccess
		title = "Quick Access"
	}

	m.telescope.SetList(items)
	m.telescope.SetTitle(title)
	m.telescope.SetTeleType(msg.teleType)

	return m, sendMsg(addStackMsg{state: STATE_TELESCOPE})
}

func (m *State) SelectTelescopeItem(msg telescope.SubmitMsg) (tea.Model, tea.Cmd) {
	val := msg.Value.Value()
	switch msg.TeleType {
	case TELESCOPE_QUICK_ACCESS:
		val, ok := val.([]tea.Msg)
		if !ok {
			return m, sendMsg(errMsg(errors.New("no menu selected")))
		}

		var cmds []tea.Cmd
		for _, msg := range val {
			cmds = append(cmds, sendMsg(msg))
		}

		return m, tea.Sequence(sendMsg(popStackMsg{}), tea.Batch(cmds...))
	case TELESCOPE_METHOD_PALLETE:
		val, ok := val.(string)
		if !ok {
			return m, sendMsg(errMsg(errors.New("no method selected")))
		}
		m.method = val
		m.url.Prompt = val + " | "
		m.url.Width = m.sw - 5 - len(m.url.Prompt)
	}

	return m, sendMsg(popStackMsg{})
}

func (m *State) CopyToClipboard() (tea.Model, tea.Cmd) {
	var val string
	switch m.state {
	case STATE_FOCUS_URL:
		val = m.url.Value()
	case STATE_FOCUS_PIPE:
		val = m.pipe.Value()
	case STATE_FOCUS_PIPEDRESP:
		val = m.pipedresp.Value()
	}

	err := clipboard.WriteAll(val)
	if err != nil {
		return m, tea.Batch(
			sendMsg(errMsg(err)),
			sendMsg(setActivityMsg("Error copying to clipboard")),
		)
	}
	return m, sendMsg(setActivityMsg("Copied to clipboard"))
}

func (m *State) Setup() (tea.Model, tea.Cmd) {
	// check whether the config folder exists
	if _, err := os.Stat(configFolder); os.IsNotExist(err) {
		err = os.MkdirAll(configFolder, 0o755)
		if err != nil {
			return m, sendMsg(errMsg(err))
		}
	}

	var saveCmd tea.Cmd
	if _, err := os.Stat(defaultSessionPath); os.IsNotExist(err) {
		saveCmd = sendMsg(saveSessionMsg{path: defaultSessionPath})
	}

	return m, tea.Batch(
		saveCmd,
		sendMsg(checkEnvFileMsg{}),
		sendMsg(loadSessionMsg{path: defaultSessionPath}),
		sendMsg(refreshStateMsg{}),
		listenResponseCmd(m.resSub),
		listenPipeResponseCmd(m.pipeResSub),
	)
}

func (m *State) HandleListFilter(msg list.FilterMatchesMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case STATE_COMMAND_PALLETE:
		m.commands, cmd = m.commands.Update(msg)
	case STATE_METHOD_PALLETE:
		m.methodSelect, cmd = m.methodSelect.Update(msg)
	case STATE_SELECT_ENV:
		m.envList, cmd = m.envList.Update(msg)
	case STATE_SELECT_SESSION, STATE_SAVE_SESSION:
		m.sessionList, cmd = m.sessionList.Update(msg)
	}

	return m, cmd
}

func (m *State) HandleRequestResult(msg requestResultMsg) (tea.Model, tea.Cmd) {
	m.response.SetValue(strings.TrimSpace(msg.res))
	return m, tea.Batch(
		sendMsg(hideSpinnerMsg{}),
		sendMsg(setActivityMsg("Request complete")),
		sendMsg(runPipeMsg{}),
		sendMsg(saveSessionMsg{path: defaultSessionPath}),
		listenResponseCmd(m.resSub),
	)
}

func (m *State) HandlePipeResult(msg pipeResultMsg) (tea.Model, tea.Cmd) {
	m.pipedresp.SetValue(strings.TrimSpace(msg.res))
	return m, tea.Batch(
		sendMsg(hideSpinnerMsg{}),
		sendMsg(setActivityMsg("Piping complete")),
		sendMsg(saveSessionMsg{path: defaultSessionPath}),
		listenPipeResponseCmd(m.pipeResSub),
	)
}

func (m *State) RunRequest() (tea.Model, tea.Cmd) {
	return m, tea.Batch(
		sendMsg(showSpinnerMsg{}),
		sendMsg(setActivityMsg("Requesting...")),
		m.Request,
	)
}

func (m *State) Request() tea.Msg {
	uc := make(chan parseWithEnvMsg)
	bc := make(chan parseWithEnvMsg)
	hc := make(chan parseWithEnvMsg)
	defer close(uc)
	defer close(bc)
	defer close(hc)

	go parseWithEnv(m.url.Value(), uc)
	go parseWithEnv(m.body.Value(), bc)
	go parseWithEnv(m.header.Value(), hc)

	url := <-uc
	if url.err != nil {
		m.resSub <- requestResultMsg{
			err: url.err,
			res: fmt.Sprintf("URL parse error : %s", url.err.Error()),
		}
		return nil
	}

	body := <-bc
	if body.err != nil {
		m.resSub <- requestResultMsg{
			err: body.err,
			res: fmt.Sprintf("Body parse error : %s", body.err.Error()),
		}
		return nil
	}

	header := <-hc
	if header.err != nil {
		m.resSub <- requestResultMsg{
			err: header.err,
			res: fmt.Sprintf("Header parse error : %s", header.err.Error()),
		}
		return nil
	}

	reqBody := requestBody{raw: body.str}
	bodyReader, err := reqBody.Buffer()
	if err != nil {
		m.resSub <- requestResultMsg{
			err: err,
			res: fmt.Sprintf("Invalid body : %s", err.Error()),
		}
		return nil
	}

	if reqBody.form != nil {
		header.str = fmt.Sprintf(
			"Content-Type: %s\n%s",
			reqBody.form.FormDataContentType(),
			header.str,
		)
	}

	req, err := http.NewRequest(m.method, url.str, bodyReader)
	if err != nil {
		m.resSub <- requestResultMsg{
			err: err,
			res: fmt.Sprintf("Request error : %s", err.Error()),
		}
		return nil
	}

	// The header example will look like this:
	// Header-1: value1
	// Header-2: value2
	reqHeader := requestHeader{raw: header.str}
	reqHeader.Apply(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		m.resSub <- requestResultMsg{
			err: err,
			res: fmt.Sprintf("Request do error : %s", err.Error()),
		}
		return nil
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		m.resSub <- requestResultMsg{
			err: err,
			res: fmt.Sprintf("Response dump error : %s", err.Error()),
		}
		return nil
	}

	var resBuffer bytes.Buffer
	resBuffer.WriteString(header.str + "\n")
	resBuffer.WriteString(responseSeparator + "\n")
	resBuffer.WriteString(body.str + "\n")
	resBuffer.WriteString(responseSeparator + "\n")
	for k, v := range resp.Header {
		resBuffer.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}
	resBuffer.WriteString(responseSeparator + "\n")
	resBuffer.WriteString(string(bodyBytes) + "\n")

	m.resSub <- requestResultMsg{res: resBuffer.String()}
	return nil
}

func (m *State) RunPipe() (tea.Model, tea.Cmd) {
	return m, tea.Batch(
		sendMsg(showSpinnerMsg{}),
		sendMsg(setActivityMsg("Piping...")),
		m.PipeRequest,
	)
}

func (m *State) PipeRequest() tea.Msg {
	resp, pipe := m.response.Value(), m.pipe.Value()
	if resp == "" {
		m.pipeResSub <- pipeResultMsg{res: resp}
		return nil
	}

	resp = m.resFilter.Filter(resp)
	if pipe == "" {
		m.pipeResSub <- pipeResultMsg{res: resp}
		return nil
	}

	command := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(
			"set -a && source %s && set +a && echo '%s' | %s",
			EnvFilePath,
			strings.ReplaceAll(resp, "'", "'\\''"),
			pipe,
		),
	)
	output, err := command.CombinedOutput()
	if err != nil {
		m.pipeResSub <- pipeResultMsg{err: err, res: string(output)}
		return nil
	}

	m.pipeResSub <- pipeResultMsg{res: string(output)}
	return nil
}

func (m *State) GetField(state string) any {
	switch state {
	case STATE_FOCUS_URL:
		return &m.url
	case STATE_FOCUS_PIPE:
		return &m.pipe
	case STATE_FOCUS_PIPEDRESP:
		return &m.pipedresp
	case STATE_FOCUS_RESPONSE:
		return &m.response
	case STATE_FOCUS_BODY:
		return &m.body
	case STATE_FOCUS_HEADER:
		return &m.header
	case STATE_COMMAND_PALLETE:
		return &m.commands
	case STATE_FOCUS_RESPONSE_FILTER:
		return &m.resFilter
	case STATE_SAVE_SESSION_INPUT, STATE_SESSION_RENAME_INPUT:
		return &m.saveInput
	}

	return nil
}

func (m *State) GetFocusedField() any {
	return m.GetField(m.state)
}

func (m *State) HandleErrorMsg(msg errMsg) (tea.Model, tea.Cmd) {
	errLog.Error(msg.Error())
	return m, nil
}

func (m *State) HandleWindowChange(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	w, h := appStyle.GetFrameSize()
	cw := msg.Width - w
	ch := msg.Height - h

	m.sw = cw
	m.sh = ch

	return m, sendMsg(recalculateComponentSizesMsg{})
}

func (m *State) NextSection() (tea.Model, tea.Cmd) {
	if !slices.Contains(homeLayout, m.state) {
		return m, nil
	}

	var nextState string
	for i, state := range homeLayout {
		if state != m.state {
			continue
		}

		nextState = homeLayout[min(len(homeLayout)-1, i+1)]
		break
	}

	return m, sendMsg(setStateMsg{state: nextState})
}

func (m *State) PrevSection() (tea.Model, tea.Cmd) {
	if !slices.Contains(homeLayout, m.state) {
		return m, nil
	}

	var prevState string
	for i, state := range homeLayout {
		if state != m.state {
			continue
		}

		prevState = homeLayout[max(0, i-1)]
		break
	}

	return m, sendMsg(setStateMsg{state: prevState})
}

func (m *State) Quit() (tea.Model, tea.Cmd) {
	if m.state == STATE_COMMAND_PALLETE && m.commands.FilterState() != list.Unfiltered {
		m.commands.ResetFilter()
		return m, nil
	}

	if m.state == STATE_SELECT_ENV && m.envList.FilterState() != list.Unfiltered {
		m.envList.ResetFilter()
		return m, nil
	}

	if m.state == STATE_METHOD_PALLETE && m.methodSelect.FilterState() != list.Unfiltered {
		m.methodSelect.ResetFilter()
		return m, nil
	}

	if slices.Contains([]string{STATE_SELECT_SESSION, STATE_SAVE_SESSION}, m.state) &&
		m.sessionList.FilterState() != list.Unfiltered {
		m.sessionList.ResetFilter()
		return m, nil
	}

	if len(m.stateStack) <= 1 {
		go m.SaveSession(saveSessionMsg{path: defaultSessionPath})
		return m, tea.Quit
	}

	return m, sendMsg(popStackMsg{})
}

func (m *State) AddStack(state string) (tea.Model, tea.Cmd) {
	m.stateStack = append(m.stateStack, state)
	m.state = state

	return m, sendMsg(refreshStateMsg{})
}

func (m *State) PopStack() (tea.Model, tea.Cmd) {
	if len(m.stateStack) <= 1 {
		return m, nil
	}

	m.stateStack = m.stateStack[:len(m.stateStack)-1]
	m.state = m.stateStack[len(m.stateStack)-1]

	return m, sendMsg(refreshStateMsg{})
}

func (m *State) PopStackRoot() (tea.Model, tea.Cmd) {
	if len(m.stateStack) <= 1 {
		return m, nil
	}

	m.stateStack = []string{m.stateStack[0]}
	m.state = m.stateStack[len(m.stateStack)-1]

	return m, sendMsg(refreshStateMsg{})
}

func (m *State) SetState(state string) (tea.Model, tea.Cmd) {
	m.state = state
	m.stateStack[len(m.stateStack)-1] = state

	return m, sendMsg(refreshStateMsg{})
}

func (m *State) RefreshState() (tea.Model, tea.Cmd) {
	m.url.Blur()
	m.response.Blur()
	m.pipe.Blur()
	m.pipedresp.Blur()
	m.body.Blur()
	m.header.Blur()
	m.resFilter.Blur()
	m.saveInput.Blur()

	switch f := m.GetFocusedField().(type) {
	case *textarea.Model:
		f.Focus()
	case *textinput.Model:
		f.Focus()
	case *ResponseFilter:
		f.Focus()
	}

	return m, nil
}

func (m *State) OpenEditor(msg openEditorMsg) (tea.Model, tea.Cmd) {
	var str string
	switch f := m.GetField(msg.state).(type) {
	case *textarea.Model:
		str = f.Value()
	case *textinput.Model:
		str = f.Value()
	default:
		return m, nil
	}

	dir := filepath.Dir(tempFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0o755)
		if err != nil {
			return m, sendMsg(errMsg(err))
		}
	}

	f, err := os.Create(tempFilePath)
	if err != nil {
		return m, sendMsg(errMsg(err))
	}

	defer f.Close()
	if _, err := f.WriteString(str); err != nil {
		return m, sendMsg(errMsg(err))
	}

	editor := getDefaultEditor()

	cmd := exec.Command(editor, tempFilePath)
	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errMsg(err)
		}

		f, err := os.Open(tempFilePath)
		if err != nil {
			return errMsg(err)
		}

		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			return errMsg(err)
		}

		str := strings.TrimSpace(string(b))

		return setFieldValueMsg{state: msg.state, value: str}
	})
}

func (m *State) SetFieldValue(msg setFieldValueMsg) (tea.Model, tea.Cmd) {
	switch f := m.GetField(msg.state).(type) {
	case *textarea.Model:
		f.SetValue(msg.value)
	case *textinput.Model:
		f.SetValue(msg.value)
	}

	return m, sendMsg(setActivityMsg("Set Field Value"))
}

func (m *State) ShowSpinner() (tea.Model, tea.Cmd) {
	m.showSpinner = true
	return m, nil
}

func (m *State) HideSpinner() (tea.Model, tea.Cmd) {
	m.showSpinner = false
	return m, nil
}

func (m *State) SelectCommandPallete() (tea.Model, tea.Cmd) {
	if m.commands.FilterState() == list.Filtering {
		var cmd tea.Cmd
		m.commands, cmd = m.commands.Update(tea.KeyMsg{
			Type: tea.KeyEnter,
		})
		return m, cmd
	}

	i, ok := m.commands.SelectedItem().(commandPallete)
	if !ok {
		return m, sendMsg(errMsg(errors.New("no command selected")))
	}

	return m, sendMsg(runCommandMsg{commandId: i.commandId})
}

func (m *State) SelectMethodPallete() (tea.Model, tea.Cmd) {
	if m.methodSelect.FilterState() == list.Filtering {
		var cmd tea.Cmd
		m.methodSelect, cmd = m.methodSelect.Update(tea.KeyMsg{
			Type: tea.KeyEnter,
		})
		return m, cmd
	}

	i, ok := m.methodSelect.SelectedItem().(methodPallete)
	if !ok {
		return m, sendMsg(errMsg(errors.New("no method selected")))
	}

	m.method = i.method
	m.url.Prompt = i.method + " | "
	m.url.Width = m.sw - 5 - len(m.url.Prompt)

	return m, sendMsg(popStackMsg{})
}

func (m *State) RunCommand(command runCommandMsg) (tea.Model, tea.Cmd) {
	switch command.commandId {
	case COMMAND_OPEN_ENV:
		return m, sendMsg(openEnvMsg{})
	case COMMAND_SELECT_METHOD:
		return m, sendMsg(addStackMsg{state: STATE_METHOD_PALLETE})
	case COMMAND_OPEN_BODY:
		return m, sendMsg(openRequestBodyMsg{})
	case COMMAND_OPEN_HEADER:
		return m, sendMsg(openRequestHeaderMsg{})
	case COMMAND_SAVE_SESSION:
		return m, tea.Sequence(
			sendMsg(addStackMsg{state: STATE_SAVE_SESSION}),
			sendMsg(loadSessionListMsg{}),
		)
	case COMMAND_OPEN_SESSION_LIST:
		return m, tea.Sequence(
			sendMsg(addStackMsg{state: STATE_SELECT_SESSION}),
			sendMsg(loadSessionListMsg{}),
		)
	case COMMAND_CHANGE_ENV:
		return m, tea.Batch(
			sendMsg(refreshSelectEnvMsg{}),
			sendMsg(addStackMsg{state: STATE_SELECT_ENV}),
		)

	default:
		return m, nil
	}
}

func (m *State) OpenEnv() (tea.Model, tea.Cmd) {
	editor := getDefaultEditor()
	cmd := exec.Command(editor, EnvFilePath)

	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errMsg(err)
		}

		return nil
	})
}

func (m *State) CheckOrCreateEnvFile() (tea.Model, tea.Cmd) {
	if _, err := os.Stat(EnvFilePath); err == nil {
		return m, sendMsg(errMsg(err))
	}

	f, err := os.Create(EnvFilePath)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return m, tea.Batch(sendMsg(errMsg(err)), tea.Quit)
	}

	defer f.Close()
	return m, nil
}

func (m *State) ReplaceCurrentSession(msg replaceCurrentSessionMsg) (tea.Model, tea.Cmd) {
	path := msg.path
	_, err := openSessionFromPath(path)
	if err != nil {
		return m, tea.Batch(
			sendMsg(errMsg(err)),
			sendMsg(setActivityMsg("Can't load session on "+path)),
		)
	}

	return m, sendMsg(loadSessionMsg{path: path})
}

func (m *State) SaveSession(msg saveSessionMsg) (tea.Model, tea.Cmd) {
	session := createSessionFromState(m)
	err := session.Save(msg.path)
	if err != nil {
		return m, tea.Batch(
			sendMsg(errMsg(err)),
			sendMsg(setActivityMsg("Error saving session")),
		)
	}

	return m, sendMsg(setActivityMsg("Session saved to " + msg.path))
}

func (m *State) LoadSession(msg loadSessionMsg) (tea.Model, tea.Cmd) {
	session, err := openSessionFromPath(msg.path)
	if err != nil {
		return m, tea.Batch(
			sendMsg(errMsg(err)),
			sendMsg(setActivityMsg("Can't load session")),
		)
	}

	session.Apply(m)

	return m, tea.Batch(
		sendMsg(setActivityMsg("Session loaded from "+msg.path)),
		sendMsg(recalculateComponentSizesMsg{}),
	)
}

func (m *State) LoadSessionList() (tea.Model, tea.Cmd) {
	if _, err := os.Stat(collectionFolder); os.IsNotExist(err) {
		return m, sendMsg(errMsg(err))
	}

	files, err := os.ReadDir(collectionFolder)
	if err != nil {
		return m, sendMsg(errMsg(err))
	}

	var sessionItems []list.Item
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		session, err := openSessionFromPath(fmt.Sprintf("%s/%s", collectionFolder, file.Name()))
		if err != nil {
			continue
		}

		sessionItems = append(sessionItems, SessionItem{name: file.Name(), session: session})
	}

	m.sessionList.SetItems(sessionItems)
	if m.sessionList.FilterState() == list.FilterApplied {
		m.sessionList.SetFilterText(m.sessionList.FilterValue())
	}

	switch m.state {
	case STATE_SAVE_SESSION:
		m.sessionList.Title = "Save Session To"
	case STATE_SELECT_SESSION:
		m.sessionList.Title = "Open Session"
	}

	return m, nil
}

func (m *State) OpenRequestBody() (tea.Model, tea.Cmd) {
	return m, sendMsg(openEditorMsg{state: STATE_FOCUS_BODY})
}

func (m *State) OpenRequestHeader() (tea.Model, tea.Cmd) {
	return m, sendMsg(openEditorMsg{state: STATE_FOCUS_HEADER})
}

func (m *State) SelectEnv() (tea.Model, tea.Cmd) {
	if m.envList.FilterState() == list.Filtering {
		var cmd tea.Cmd
		m.envList, cmd = m.envList.Update(tea.KeyMsg{
			Type: tea.KeyEnter,
		})
		return m, cmd
	}

	file, ok := m.envList.SelectedItem().(fileItem)
	if !ok {
		return m, sendMsg(errMsg(errors.New("no env selected")))
	}

	EnvFilePath = file.name

	return m, sendMsg(popStackMsg{})
}

func (m *State) RefreshSelectEnv() (tea.Model, tea.Cmd) {
	files, err := os.ReadDir(".")
	if err != nil {
		return m, sendMsg(errMsg(err))
	}

	var fileItems []list.Item
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileItems = append(fileItems, fileItem{name: file.Name(), path: file.Name()})
	}

	m.envList.SetItems(fileItems)
	m.envList.KeyMap.Quit.SetEnabled(false)
	return m, nil
}

func (m *State) DeleteSessionItem() (tea.Model, tea.Cmd) {
	i, ok := m.sessionList.SelectedItem().(SessionItem)
	if !ok {
		return m, sendMsg(errMsg(errors.New("no session selected")))
	}

	err := i.Delete()
	if err != nil {
		return m, sendMsg(errMsg(err))
	}

	return m, sendMsg(loadSessionListMsg{})
}

func (m *State) SelectSessionItem() (tea.Model, tea.Cmd) {
	if m.sessionList.FilterState() == list.Filtering {
		var cmd tea.Cmd
		m.sessionList, cmd = m.sessionList.Update(tea.KeyMsg{
			Type: tea.KeyEnter,
		})
		return m, cmd
	}

	i, ok := m.sessionList.SelectedItem().(SessionItem)
	if !ok {
		return m, sendMsg(errMsg(errors.New("no session selected")))
	}

	switch m.state {
	case STATE_SAVE_SESSION:
		return m, tea.Batch(sendMsg(saveSessionMsg{path: i.Path()}), sendMsg(popStackRootMsg{}))
	case STATE_SELECT_SESSION:
		return m, tea.Batch(
			sendMsg(replaceCurrentSessionMsg{path: i.Path()}),
			sendMsg(popStackRootMsg{}),
		)
	}

	return m, nil
}

func (m *State) SaveSessionInputSubmit() (tea.Model, tea.Cmd) {
	val := strings.TrimSpace(m.saveInput.Value())
	if val == "" {
		return m, sendMsg(errMsg(errors.New("new name is empty")))
	}

	switch m.state {
	case STATE_SAVE_SESSION_INPUT:
		return m, tea.Batch(
			sendMsg(saveSessionMsg{path: fmt.Sprintf("%s/%s.json", collectionFolder, val)}),
			sendMsg(popStackRootMsg{}),
		)
	case STATE_SESSION_RENAME_INPUT:
		i, ok := m.sessionList.SelectedItem().(SessionItem)
		if !ok {
			return m, sendMsg(errMsg(errors.New("no session selected")))
		}

		err := i.Rename(val)
		if err != nil {
			return m, sendMsg(errMsg(err))
		}

		return m, tea.Batch(
			sendMsg(popStackMsg{}),
			sendMsg(loadSessionListMsg{}),
		)
	}

	return m, nil
}
