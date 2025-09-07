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

	c "github.com/Cybernetics354/mayohttp/internal/app/common"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/event"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/intf"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/store"
	"github.com/Cybernetics354/mayohttp/internal/app/ui/telescope"
	"github.com/Cybernetics354/mayohttp/internal/app/val"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *State) applySession(s *store.Session) *State {
	m.url.SetValue(s.Url)
	m.pipe.SetValue(s.Pipe)
	m.pipedresp.SetValue(s.PipedResponse)
	m.response.SetValue(s.Response)
	m.header.SetValue(s.Header)
	m.body.SetValue(s.Body)
	m.method = s.Method
	m.resFilter = s.ResFilter
	m.url.Prompt = m.method + " | "
	m.url.Width = m.sw - 5 - len(m.url.Prompt)

	return m
}

func createSessionFromState(s *State) *store.Session {
	return &store.Session{
		Url:           s.url.Value(),
		Pipe:          s.pipe.Value(),
		PipedResponse: s.pipedresp.Value(),
		Method:        s.method,
		Response:      s.response.Value(),
		Header:        s.header.Value(),
		Body:          s.body.Value(),
		ResFilter:     s.resFilter,
	}
}

func (m *State) ClearFocusedInput() (tea.Model, tea.Cmd) {
	field := m.GetFocusedField()
	if field == nil {
		return m, nil
	}

	switch f := field.(type) {
	case *textinput.Model:
		f.SetValue("")
	}

	return m, nil
}

func (m *State) OpenTelescope(msg event.OpenTelescope) (tea.Model, tea.Cmd) {
	var items []list.Item
	var title string

	switch msg.TeleType {
	case val.TELESCOPE_METHOD_PALLETE:
		items = val.MethodPalletesTelescope
		title = "Select Method"
	case val.TELESCOPE_QUICK_ACCESS:
		items = val.QuickAccess
		title = "Quick Access"
	}

	m.telescope.SetList(items)
	m.telescope.SetTitle(title)
	m.telescope.SetTeleType(msg.TeleType)

	return m, c.SendMsg(event.AddStack{State: val.STATE_TELESCOPE})
}

func (m *State) SelectTelescopeItem(msg telescope.SubmitMsg) (tea.Model, tea.Cmd) {
	switch msg.TeleType {
	case val.TELESCOPE_QUICK_ACCESS:
		v, ok := msg.Value.Value().([]tea.Msg)
		if !ok {
			return m, c.SendMsg(event.Err(errors.New("no menu selected")))
		}

		var cmds []tea.Cmd
		for _, msg := range v {
			cmds = append(cmds, c.SendMsg(msg))
		}

		return m, tea.Sequence(c.SendMsg(event.PopStack{}), tea.Batch(cmds...))
	case val.TELESCOPE_METHOD_PALLETE:
		v, ok := msg.Value.Value().(string)
		if !ok {
			return m, c.SendMsg(event.Err(errors.New("no method selected")))
		}

		m.telescope.Clear()
		m.method = v
		m.url.Prompt = v + " | "
		m.url.Width = m.sw - 5 - len(m.url.Prompt)
	}

	return m, c.SendMsg(event.PopStack{})
}

func (m *State) CopyToClipboard() (tea.Model, tea.Cmd) {
	var v string
	switch m.state {
	case val.STATE_FOCUS_URL:
		v = m.url.Value()
	case val.STATE_FOCUS_PIPE:
		v = m.pipe.Value()
	case val.STATE_FOCUS_PIPEDRESP:
		v = m.pipedresp.Value()
	}

	err := clipboard.WriteAll(v)
	if err != nil {
		return m, tea.Batch(
			c.SendMsg(event.Err(err)),
			c.SendMsg(event.SetActivity("Error copying to clipboard")),
		)
	}
	return m, c.SendMsg(event.SetActivity("Copied to clipboard"))
}

func (m *State) Setup() (tea.Model, tea.Cmd) {
	// check whether the config folder exists
	if _, err := os.Stat(val.ConfigFolder); os.IsNotExist(err) {
		err = os.MkdirAll(val.ConfigFolder, 0o755)
		if err != nil {
			return m, c.SendMsg(event.Err(err))
		}
	}

	var saveCmd tea.Cmd
	if _, err := os.Stat(val.DefaultSessionPath); os.IsNotExist(err) {
		saveCmd = c.SendMsg(event.SaveSession{Path: val.DefaultSessionPath})
	}

	return m, tea.Batch(
		saveCmd,
		c.SendMsg(event.CheckEnvFile{}),
		c.SendMsg(event.LoadSession{Path: val.DefaultSessionPath}),
		c.SendMsg(event.RefreshState{}),
		event.ListenResponseCmd(m.resSub),
		event.ListenPipeResponseCmd(m.pipeResSub),
	)
}

func (m *State) HandleListFilter(msg list.FilterMatchesMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case val.STATE_COMMAND_PALLETE:
		m.commands, cmd = m.commands.Update(msg)
	case val.STATE_METHOD_PALLETE:
		m.methodSelect, cmd = m.methodSelect.Update(msg)
	case val.STATE_SELECT_ENV:
		m.envList, cmd = m.envList.Update(msg)
	case val.STATE_SELECT_SESSION, val.STATE_SAVE_SESSION:
		m.sessionList, cmd = m.sessionList.Update(msg)
	}

	return m, cmd
}

func (m *State) HandleRequestResult(msg event.RequestResult) (tea.Model, tea.Cmd) {
	m.response.SetValue(strings.TrimSpace(msg.Res))
	return m, tea.Batch(
		c.SendMsg(event.HideSpinner{}),
		c.SendMsg(event.SetActivity("Request complete")),
		c.SendMsg(event.RunPipe{}),
		c.SendMsg(event.SaveSession{Path: val.DefaultSessionPath}),
		event.ListenResponseCmd(m.resSub),
	)
}

func (m *State) HandlePipeResult(msg event.PipeResult) (tea.Model, tea.Cmd) {
	m.pipedresp.SetValue(strings.TrimSpace(msg.Res))
	return m, tea.Batch(
		c.SendMsg(event.HideSpinner{}),
		c.SendMsg(event.SetActivity("Piping complete")),
		c.SendMsg(event.SaveSession{Path: val.DefaultSessionPath}),
		event.ListenPipeResponseCmd(m.pipeResSub),
	)
}

func (m *State) RunRequest() (tea.Model, tea.Cmd) {
	return m, tea.Batch(
		c.SendMsg(event.ShowSpinner{}),
		c.SendMsg(event.SetActivity("Requesting...")),
		m.Request,
	)
}

func (m *State) Request() tea.Msg {
	uc := make(chan c.ParseWithEnvMsg)
	bc := make(chan c.ParseWithEnvMsg)
	hc := make(chan c.ParseWithEnvMsg)
	defer close(uc)
	defer close(bc)
	defer close(hc)

	go c.ParseWithEnv(m.url.Value(), uc)
	go c.ParseWithEnv(m.body.Value(), bc)
	go c.ParseWithEnv(m.header.Value(), hc)

	url := <-uc
	if url.Err != nil {
		m.resSub <- event.RequestResult{
			Err: url.Err,
			Res: fmt.Sprintf("URL parse error : %s", url.Err.Error()),
		}
		return nil
	}

	body := <-bc
	if body.Err != nil {
		m.resSub <- event.RequestResult{
			Err: body.Err,
			Res: fmt.Sprintf("Body parse error : %s", body.Err.Error()),
		}
		return nil
	}

	header := <-hc
	if header.Err != nil {
		m.resSub <- event.RequestResult{
			Err: header.Err,
			Res: fmt.Sprintf("Header parse error : %s", header.Err.Error()),
		}
		return nil
	}

	reqBody := intf.RequestBody{Raw: body.Str}
	bodyReader, err := reqBody.Buffer()
	if err != nil {
		m.resSub <- event.RequestResult{
			Err: err,
			Res: fmt.Sprintf("Invalid body : %s", err.Error()),
		}
		return nil
	}

	if reqBody.Form != nil {
		header.Str = fmt.Sprintf(
			"Content-Type: %s\n%s",
			reqBody.Form.FormDataContentType(),
			header.Str,
		)
	}

	req, err := http.NewRequest(m.method, url.Str, bodyReader)
	if err != nil {
		m.resSub <- event.RequestResult{
			Err: err,
			Res: fmt.Sprintf("Request error : %s", err.Error()),
		}
		return nil
	}

	// The header example will look like this:
	// Header-1: value1
	// Header-2: value2
	reqHeader := intf.RequestHeader{Raw: header.Str}
	reqHeader.Apply(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		m.resSub <- event.RequestResult{
			Err: err,
			Res: fmt.Sprintf("Request do error : %s", err.Error()),
		}
		return nil
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		m.resSub <- event.RequestResult{
			Err: err,
			Res: fmt.Sprintf("Response dump error : %s", err.Error()),
		}
		return nil
	}

	var resBuffer bytes.Buffer
	resBuffer.WriteString(header.Str + "\n")
	resBuffer.WriteString(val.ResponseSeparator + "\n")
	resBuffer.WriteString(body.Str + "\n")
	resBuffer.WriteString(val.ResponseSeparator + "\n")
	for k, v := range resp.Header {
		resBuffer.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}
	resBuffer.WriteString(val.ResponseSeparator + "\n")
	resBuffer.WriteString(string(bodyBytes) + "\n")

	m.resSub <- event.RequestResult{Res: resBuffer.String()}
	return nil
}

func (m *State) RunPipe() (tea.Model, tea.Cmd) {
	return m, tea.Batch(
		c.SendMsg(event.ShowSpinner{}),
		c.SendMsg(event.SetActivity("Piping...")),
		m.PipeRequest,
	)
}

func (m *State) PipeRequest() tea.Msg {
	resp, pipe := m.response.Value(), m.pipe.Value()
	if resp == "" {
		m.pipeResSub <- event.PipeResult{Res: resp}
		return nil
	}

	resp = m.resFilter.Filter(resp)
	if pipe == "" {
		m.pipeResSub <- event.PipeResult{Res: resp}
		return nil
	}

	command := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(
			"set -a && source %s && set +a && echo '%s' | %s",
			val.EnvFilePath,
			strings.ReplaceAll(resp, "'", "'\\''"),
			pipe,
		),
	)
	output, err := command.CombinedOutput()
	if err != nil {
		m.pipeResSub <- event.PipeResult{Err: err, Res: string(output)}
		return nil
	}

	m.pipeResSub <- event.PipeResult{Res: string(output)}
	return nil
}

func (m *State) GetField(state string) any {
	switch state {
	case val.STATE_FOCUS_URL:
		return &m.url
	case val.STATE_FOCUS_PIPE:
		return &m.pipe
	case val.STATE_FOCUS_PIPEDRESP:
		return &m.pipedresp
	case val.STATE_FOCUS_RESPONSE:
		return &m.response
	case val.STATE_FOCUS_BODY:
		return &m.body
	case val.STATE_FOCUS_HEADER:
		return &m.header
	case val.STATE_COMMAND_PALLETE:
		return &m.commands
	case val.STATE_FOCUS_RESPONSE_FILTER:
		return &m.resFilter
	case val.STATE_SAVE_SESSION_INPUT, val.STATE_SESSION_RENAME_INPUT:
		return &m.saveInput
	}

	return nil
}

func (m *State) GetFocusedField() any {
	return m.GetField(m.state)
}

func (m *State) HandleErrorMsg(msg event.Err) (tea.Model, tea.Cmd) {
	c.ErrLog.Error(msg.Error())
	return m, nil
}

func (m *State) HandleWindowChange(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	w, h := val.AppStyle.GetFrameSize()
	cw := msg.Width - w
	ch := msg.Height - h

	m.sw = cw
	m.sh = ch

	return m, c.SendMsg(event.RecalculateComponentSizes{})
}

func (m *State) NextSection() (tea.Model, tea.Cmd) {
	if !slices.Contains(val.HomeLayout, m.state) {
		return m, nil
	}

	var nextState string
	for i, state := range val.HomeLayout {
		if state != m.state {
			continue
		}

		nextState = val.HomeLayout[min(len(val.HomeLayout)-1, i+1)]
		break
	}

	return m, c.SendMsg(event.SetState{State: nextState})
}

func (m *State) PrevSection() (tea.Model, tea.Cmd) {
	if !slices.Contains(val.HomeLayout, m.state) {
		return m, nil
	}

	var prevState string
	for i, state := range val.HomeLayout {
		if state != m.state {
			continue
		}

		prevState = val.HomeLayout[max(0, i-1)]
		break
	}

	return m, c.SendMsg(event.SetState{State: prevState})
}

func (m *State) Quit() (tea.Model, tea.Cmd) {
	if m.state == val.STATE_COMMAND_PALLETE && m.commands.FilterState() != list.Unfiltered {
		m.commands.ResetFilter()
		return m, nil
	}

	if m.state == val.STATE_SELECT_ENV && m.envList.FilterState() != list.Unfiltered {
		m.envList.ResetFilter()
		return m, nil
	}

	if m.state == val.STATE_METHOD_PALLETE && m.methodSelect.FilterState() != list.Unfiltered {
		m.methodSelect.ResetFilter()
		return m, nil
	}

	if slices.Contains([]string{val.STATE_SELECT_SESSION, val.STATE_SAVE_SESSION}, m.state) &&
		m.sessionList.FilterState() != list.Unfiltered {
		m.sessionList.ResetFilter()
		return m, nil
	}

	if len(m.stateStack) <= 1 {
		go m.SaveSession(event.SaveSession{Path: val.DefaultSessionPath})
		return m, tea.Quit
	}

	return m, c.SendMsg(event.PopStack{})
}

func (m *State) AddStack(state string) (tea.Model, tea.Cmd) {
	m.stateStack = append(m.stateStack, state)
	m.state = state

	return m, c.SendMsg(event.RefreshState{})
}

func (m *State) PopStack() (tea.Model, tea.Cmd) {
	if len(m.stateStack) <= 1 {
		return m, nil
	}

	m.stateStack = m.stateStack[:len(m.stateStack)-1]
	m.state = m.stateStack[len(m.stateStack)-1]

	return m, c.SendMsg(event.RefreshState{})
}

func (m *State) PopStackRoot() (tea.Model, tea.Cmd) {
	if len(m.stateStack) <= 1 {
		return m, nil
	}

	m.stateStack = []string{m.stateStack[0]}
	m.state = m.stateStack[len(m.stateStack)-1]

	return m, c.SendMsg(event.RefreshState{})
}

func (m *State) SetState(state string) (tea.Model, tea.Cmd) {
	m.state = state
	m.stateStack[len(m.stateStack)-1] = state

	return m, c.SendMsg(event.RefreshState{})
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
	case *intf.ResponseFilter:
		f.Focus()
	}

	return m, nil
}

func (m *State) OpenEditor(msg event.OpenEditor) (tea.Model, tea.Cmd) {
	var str string
	switch f := m.GetField(msg.State).(type) {
	case *textarea.Model:
		str = f.Value()
	case *textinput.Model:
		str = f.Value()
	default:
		return m, nil
	}

	dir := filepath.Dir(val.TempFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0o755)
		if err != nil {
			return m, c.SendMsg(event.Err(err))
		}
	}

	f, err := os.Create(val.TempFilePath)
	if err != nil {
		return m, c.SendMsg(event.Err(err))
	}

	defer f.Close()
	if _, err := f.WriteString(str); err != nil {
		return m, c.SendMsg(event.Err(err))
	}

	editor := c.GetDefaultEditor()

	cmd := exec.Command(editor, val.TempFilePath)
	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return event.Err(err)
		}

		f, err := os.Open(val.TempFilePath)
		if err != nil {
			return event.Err(err)
		}

		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			return event.Err(err)
		}

		str := strings.TrimSpace(string(b))

		return event.SetFieldValue{State: msg.State, Value: str}
	})
}

func (m *State) SetFieldValue(msg event.SetFieldValue) (tea.Model, tea.Cmd) {
	switch f := m.GetField(msg.State).(type) {
	case *textarea.Model:
		f.SetValue(msg.Value)
	case *textinput.Model:
		f.SetValue(msg.Value)
	}

	return m, c.SendMsg(event.SetActivity("Set Field Value"))
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

	i, ok := m.commands.SelectedItem().(val.CommandPallete)
	if !ok {
		return m, c.SendMsg(event.Err(errors.New("no command selected")))
	}

	return m, c.SendMsg(event.RunCommand{CommandId: i.CommandId})
}

func (m *State) SelectMethodPallete() (tea.Model, tea.Cmd) {
	if m.methodSelect.FilterState() == list.Filtering {
		var cmd tea.Cmd
		m.methodSelect, cmd = m.methodSelect.Update(tea.KeyMsg{
			Type: tea.KeyEnter,
		})
		return m, cmd
	}

	i, ok := m.methodSelect.SelectedItem().(val.MethodPallete)
	if !ok {
		return m, c.SendMsg(event.Err(errors.New("no method selected")))
	}

	m.method = i.Method
	m.url.Prompt = i.Method + " | "
	m.url.Width = m.sw - 5 - len(m.url.Prompt)

	return m, c.SendMsg(event.PopStack{})
}

func (m *State) RunCommand(command event.RunCommand) (tea.Model, tea.Cmd) {
	switch command.CommandId {
	case val.COMMAND_OPEN_ENV:
		return m, c.SendMsg(event.OpenEnv{})
	case val.COMMAND_SELECT_METHOD:
		return m, c.SendMsg(event.OpenTelescope{TeleType: val.TELESCOPE_METHOD_PALLETE})
	case val.COMMAND_OPEN_BODY:
		return m, c.SendMsg(event.OpenRequestBody{})
	case val.COMMAND_OPEN_HEADER:
		return m, c.SendMsg(event.OpenRequestHeader{})
	case val.COMMAND_SAVE_SESSION:
		return m, tea.Sequence(
			c.SendMsg(event.AddStack{State: val.STATE_SAVE_SESSION}),
			c.SendMsg(event.LoadSessionList{}),
		)
	case val.COMMAND_OPEN_SESSION_LIST:
		return m, tea.Sequence(
			c.SendMsg(event.AddStack{State: val.STATE_SELECT_SESSION}),
			c.SendMsg(event.LoadSessionList{}),
		)
	case val.COMMAND_CHANGE_ENV:
		return m, tea.Batch(
			c.SendMsg(event.RefreshSelectEnv{}),
			c.SendMsg(event.AddStack{State: val.STATE_SELECT_ENV}),
		)

	default:
		return m, nil
	}
}

func (m *State) OpenEnv() (tea.Model, tea.Cmd) {
	editor := c.GetDefaultEditor()
	cmd := exec.Command(editor, val.EnvFilePath)

	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return event.Err(err)
		}

		return nil
	})
}

func (m *State) CheckOrCreateEnvFile() (tea.Model, tea.Cmd) {
	if _, err := os.Stat(val.EnvFilePath); err == nil {
		return m, c.SendMsg(event.Err(err))
	}

	f, err := os.Create(val.EnvFilePath)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return m, tea.Batch(c.SendMsg(event.Err(err)), tea.Quit)
	}

	defer f.Close()
	return m, nil
}

func (m *State) ReplaceCurrentSession(msg event.ReplaceCurrentSession) (tea.Model, tea.Cmd) {
	path := msg.Path
	_, err := store.OpenSessionFromPath(path)
	if err != nil {
		return m, tea.Batch(
			c.SendMsg(event.Err(err)),
			c.SendMsg(event.SetActivity("Can't load session on "+path)),
		)
	}

	return m, c.SendMsg(event.LoadSession{Path: path})
}

func (m *State) SaveSession(msg event.SaveSession) (tea.Model, tea.Cmd) {
	session := createSessionFromState(m)
	err := session.Save(msg.Path)
	if err != nil {
		return m, tea.Batch(
			c.SendMsg(event.Err(err)),
			c.SendMsg(event.SetActivity("Error saving session")),
		)
	}

	return m, c.SendMsg(event.SetActivity("Session saved to " + msg.Path))
}

func (m *State) LoadSession(msg event.LoadSession) (tea.Model, tea.Cmd) {
	session, err := store.OpenSessionFromPath(msg.Path)
	if err != nil {
		return m, tea.Batch(
			c.SendMsg(event.Err(err)),
			c.SendMsg(event.SetActivity("Can't load session")),
		)
	}

	m.applySession(session)

	return m, tea.Batch(
		c.SendMsg(event.SetActivity("Session loaded from "+msg.Path)),
		c.SendMsg(event.RecalculateComponentSizes{}),
	)
}

func (m *State) LoadSessionList() (tea.Model, tea.Cmd) {
	if _, err := os.Stat(val.CollectionFolder); os.IsNotExist(err) {
		return m, c.SendMsg(event.Err(err))
	}

	files, err := os.ReadDir(val.CollectionFolder)
	if err != nil {
		return m, c.SendMsg(event.Err(err))
	}

	var sessionItems []list.Item
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		session, err := store.OpenSessionFromPath(
			fmt.Sprintf("%s/%s", val.CollectionFolder, file.Name()),
		)
		if err != nil {
			continue
		}

		session.Name = file.Name()

		sessionItems = append(sessionItems, session)
	}

	m.sessionList.SetItems(sessionItems)
	if m.sessionList.FilterState() == list.FilterApplied {
		m.sessionList.SetFilterText(m.sessionList.FilterValue())
	}

	switch m.state {
	case val.STATE_SAVE_SESSION:
		m.sessionList.Title = "Save Session To"
	case val.STATE_SELECT_SESSION:
		m.sessionList.Title = "Open Session"
	}

	return m, nil
}

func (m *State) OpenRequestBody() (tea.Model, tea.Cmd) {
	return m, c.SendMsg(event.OpenEditor{State: val.STATE_FOCUS_BODY})
}

func (m *State) OpenRequestHeader() (tea.Model, tea.Cmd) {
	return m, c.SendMsg(event.OpenEditor{State: val.STATE_FOCUS_HEADER})
}

func (m *State) SelectEnv() (tea.Model, tea.Cmd) {
	if m.envList.FilterState() == list.Filtering {
		var cmd tea.Cmd
		m.envList, cmd = m.envList.Update(tea.KeyMsg{
			Type: tea.KeyEnter,
		})
		return m, cmd
	}

	file, ok := m.envList.SelectedItem().(intf.FileItem)
	if !ok {
		return m, c.SendMsg(event.Err(errors.New("no env selected")))
	}

	val.EnvFilePath = file.Name

	return m, c.SendMsg(event.PopStack{})
}

func (m *State) RefreshSelectEnv() (tea.Model, tea.Cmd) {
	files, err := os.ReadDir(".")
	if err != nil {
		return m, c.SendMsg(event.Err(err))
	}

	var fileItems []list.Item
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileItems = append(fileItems, intf.FileItem{Name: file.Name(), Path: file.Name()})
	}

	m.envList.SetItems(fileItems)
	m.envList.KeyMap.Quit.SetEnabled(false)
	return m, nil
}

func (m *State) DeleteSessionItem() (tea.Model, tea.Cmd) {
	i, ok := m.sessionList.SelectedItem().(*store.Session)
	if !ok {
		return m, c.SendMsg(event.Err(errors.New("no session selected")))
	}

	err := i.Delete()
	if err != nil {
		return m, c.SendMsg(event.Err(err))
	}

	return m, c.SendMsg(event.LoadSessionList{})
}

func (m *State) SelectSessionItem() (tea.Model, tea.Cmd) {
	if m.sessionList.FilterState() == list.Filtering {
		var cmd tea.Cmd
		m.sessionList, cmd = m.sessionList.Update(tea.KeyMsg{
			Type: tea.KeyEnter,
		})
		return m, cmd
	}

	i, ok := m.sessionList.SelectedItem().(*store.Session)
	if !ok {
		return m, c.SendMsg(event.Err(errors.New("no session selected")))
	}

	switch m.state {
	case val.STATE_SAVE_SESSION:
		return m, tea.Batch(
			c.SendMsg(event.SaveSession{Path: i.Path()}),
			c.SendMsg(event.PopStackRoot{}),
		)
	case val.STATE_SELECT_SESSION:
		return m, tea.Batch(
			c.SendMsg(event.ReplaceCurrentSession{Path: i.Path()}),
			c.SendMsg(event.PopStackRoot{}),
		)
	}

	return m, nil
}

func (m *State) SaveSessionInputSubmit() (tea.Model, tea.Cmd) {
	v := strings.TrimSpace(m.saveInput.Value())
	if v == "" {
		return m, c.SendMsg(event.Err(errors.New("new name is empty")))
	}

	switch m.state {
	case val.STATE_SAVE_SESSION_INPUT:
		return m, tea.Batch(
			c.SendMsg(
				event.SaveSession{Path: fmt.Sprintf("%s/%s.json", val.CollectionFolder, v)},
			),
			c.SendMsg(event.PopStackRoot{}),
		)
	case val.STATE_SESSION_RENAME_INPUT:
		i, ok := m.sessionList.SelectedItem().(*store.Session)
		if !ok {
			return m, c.SendMsg(event.Err(errors.New("no session selected")))
		}

		err := i.Rename(v)
		if err != nil {
			return m, c.SendMsg(event.Err(err))
		}

		return m, tea.Batch(
			c.SendMsg(event.PopStack{}),
			c.SendMsg(event.LoadSessionList{}),
		)
	}

	return m, nil
}
