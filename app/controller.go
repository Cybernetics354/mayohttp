package app

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m State) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("Mayo HTTP"),
		tea.EnterAltScreen,
		m.spinner.Tick,
		refreshState,
		listenResponse(m.resSub),
		listenPipeResponse(m.pipeResSub),
	)
}

func (m State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case list.FilterMatchesMsg:
		m.commands, cmd = m.commands.Update(msg)
		return m, cmd
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.WindowSizeMsg:
		return m.HandleWindowChange(msg)
	case tea.KeyMsg:
		return m.HandleKeyPress(msg)
	case recalculateComponentSizesMsg:
		return m.RecalculateComponentSize()
	case setStateMsg:
		return m.SetState(msg.state)
	case addStackMsg:
		return m.AddStack(msg.state)
	case popStackMsg:
		return m.PopStack()
	case nextSectionMsg:
		return m.NextSection()
	case prevSectionMsg:
		return m.PrevSection()
	case selectCommandPalleteMsg:
		return m.SelectCommandPallete()
	case runCommandMsg:
		return m.RunCommand(msg)
	case openEnvMsg:
		return m.OpenEnv()
	case openEditorMsg:
		return m.OpenEditor()
	case hideSpinnerMsg:
		return m.HideSpinner()
	case showSpinnerMsg:
		return m.ShowSpinner()
	case refreshStateMsg:
		return m.RefreshState()
	case runRequestMsg:
		return m.RunRequest()
	case runPipeMsg:
		return m.RunPipe()
	case requestResultMsg:
		return m.HandleRequestResult(msg)
	case pipeResultMsg:
		return m.HandlePipeResult(msg)
	case errMsg:
		return m.HandleErrorMsg(msg)
	}

	return m, nil
}

func (m State) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	return m.Render()
}

func (m *State) HandleRequestResult(msg requestResultMsg) (tea.Model, tea.Cmd) {
	m.response.SetValue(msg.res)
	return m, tea.Batch(hideSpinner, runPipe, listenResponse(m.resSub))
}

func (m *State) HandlePipeResult(msg pipeResultMsg) (tea.Model, tea.Cmd) {
	m.pipedresp.SetValue(msg.res)
	return m, tea.Batch(hideSpinner, listenPipeResponse(m.pipeResSub))
}

func (m *State) RunRequest() (tea.Model, tea.Cmd) {
	return m, tea.Batch(showSpinner, m.Request)
}

func (m *State) Request() tea.Msg {
	command := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(
			"set -a && source %s && set +a && echo '%s' | envsubst",
			EnvFilePath,
			m.url.Value(),
		),
	)

	url, err := command.Output()
	if err != nil {
		m.resSub <- requestResultMsg{
			err: err,
			res: fmt.Sprintf("Failed to parse url : %s", err.Error()),
		}
		return nil
	}

	req, err := http.NewRequest(m.method, strings.TrimSpace(string(url)), nil)
	if err != nil {
		m.resSub <- requestResultMsg{
			err: err,
			res: fmt.Sprintf("Request error : %s", err.Error()),
		}
		return nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		m.resSub <- requestResultMsg{
			err: err,
			res: fmt.Sprintf("Request do error : %s", err.Error()),
		}
		return nil
	}

	defer resp.Body.Close()

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		m.resSub <- requestResultMsg{
			err: err,
			res: fmt.Sprintf("Response dump error : %s", err.Error()),
		}
		return nil
	}

	m.resSub <- requestResultMsg{res: string(respDump)}
	return nil
}

func (m *State) RunPipe() (tea.Model, tea.Cmd) {
	return m, tea.Batch(showSpinner, m.PipeRequest)
}

func (m *State) PipeRequest() tea.Msg {
	resp, pipe := m.response.Value(), m.pipe.Value()
	if resp == "" || pipe == "" {
		m.pipeResSub <- pipeResultMsg{res: m.response.Value()}
		return nil
	}

	command := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf("set -a && source %s && set +a && echo '%s' | %s", EnvFilePath, resp, pipe),
	)
	output, err := command.CombinedOutput()
	if err != nil {
		m.pipeResSub <- pipeResultMsg{err: err, res: string(output)}
		return nil
	}

	m.pipeResSub <- pipeResultMsg{res: string(output)}
	return nil
}

func (m *State) GetFocusedField() any {
	switch m.state {
	case FOCUS_URL:
		return &m.url
	case FOCUS_PIPE:
		return &m.pipe
	case FOCUS_PIPEDRESP:
		return &m.pipedresp
	case FOCUS_RESPONSE:
		return &m.response
	case EDIT_BODY:
		return &m.body
	case EDIT_HEADER:
		return &m.header
	case COMMAND_PALLETE:
		return &m.commands
	}

	return &m.url
}

func (m *State) HandleErrorMsg(msg errMsg) (tea.Model, tea.Cmd) {
	m.err = msg
	return m, nil
}

func (m *State) HandleWindowChange(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	w, h := appStyle.GetFrameSize()
	cw := msg.Width - w
	ch := msg.Height - h

	m.sw = cw
	m.sh = ch

	return m, recalculateComponentSizes
}

func (m *State) HandleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, keyMaps.Open):
		return m, openEditor
	case key.Matches(msg, keyMaps.Quit):
		return m.Quit()
	case key.Matches(msg, keyMaps.Commands):
		return m, addStack(COMMAND_PALLETE)
	case key.Matches(msg, keyMaps.Next):
		return m, nextSection
	case key.Matches(msg, keyMaps.Back):
		return m, prevSection
	case key.Matches(msg, keyMaps.Run):
		switch m.state {
		case FOCUS_URL:
			return m, runRequest
		case FOCUS_PIPE:
			return m, runPipe
		case COMMAND_PALLETE:
			return m, selectCommandPallete
		}

		if m.state == FOCUS_URL {
			return m, runRequest
		}

		if m.state == FOCUS_PIPE {
			return m, runPipe
		}

		if m.state == COMMAND_PALLETE {
			return m, selectCommandPallete
		}
	}

	switch m.state {
	case FOCUS_PIPE:
		m.pipe, cmd = m.pipe.Update(msg)
	case FOCUS_URL:
		m.url, cmd = m.url.Update(msg)
	case FOCUS_RESPONSE:
		m.response, cmd = m.response.Update(msg)
	case FOCUS_PIPEDRESP:
		m.pipedresp, cmd = m.pipedresp.Update(msg)
	case COMMAND_PALLETE:
		m.commands, cmd = m.commands.Update(msg)
	}

	return m, cmd
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

	return m, setState(nextState)
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

	return m, setState(prevState)
}

func (m *State) Quit() (tea.Model, tea.Cmd) {
	if len(m.stateStack) <= 1 {
		return m, tea.Quit
	}

	return m, popStack
}

func (m *State) AddStack(state string) (tea.Model, tea.Cmd) {
	m.stateStack = append(m.stateStack, state)
	m.state = state

	return m, refreshState
}

func (m *State) PopStack() (tea.Model, tea.Cmd) {
	if len(m.stateStack) <= 1 {
		return m, nil
	}

	m.stateStack = m.stateStack[:len(m.stateStack)-1]
	m.state = m.stateStack[len(m.stateStack)-1]

	return m, refreshState
}

func (m *State) SetState(state string) (tea.Model, tea.Cmd) {
	m.state = state
	m.stateStack[len(m.stateStack)-1] = state

	return m, refreshState
}

func (m *State) RefreshState() (tea.Model, tea.Cmd) {
	m.url.Blur()
	m.response.Blur()
	m.pipe.Blur()
	m.pipedresp.Blur()
	m.body.Blur()
	m.header.Blur()
	m.commands.ResetFilter()

	switch f := m.GetFocusedField().(type) {
	case *textarea.Model:
		f.Focus()
	case *textinput.Model:
		f.Focus()
	}

	return m, nil
}

func (m *State) OpenEditor() (tea.Model, tea.Cmd) {
	var str string
	switch f := m.GetFocusedField().(type) {
	case *textarea.Model:
		str = f.Value()
	case *textinput.Model:
		str = f.Value()
	default:
		return m, nil
	}

	dir := filepath.Dir(tempFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			m.err = err
			return m, nil
		}
	}

	f, err := os.Create(tempFilePath)
	if err != nil {
		m.err = err
		return m, nil
	}

	f.WriteString(str)
	f.Close()

	editor := getDefaultEditor()

	cmd := exec.Command(editor, tempFilePath)
	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
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

		switch f := m.GetFocusedField().(type) {
		case *textarea.Model:
			f.SetValue(str)
		case *textinput.Model:
			f.SetValue(str)
		}

		return errMsg(err)
	})
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
	i, ok := m.commands.SelectedItem().(commandPallete)
	if !ok {
		return m, popStack
	}

	return m, tea.Batch(runCommand(i.commandId), popStack)
}

func (m *State) RunCommand(command runCommandMsg) (tea.Model, tea.Cmd) {
	switch command.commandId {
	case COMMAND_OPEN_ENV:
		return m, tea.Batch(openEnv, nil)
	default:
		return m, nil
	}
}

func (m *State) OpenEnv() (tea.Model, tea.Cmd) {
	editor := getDefaultEditor()
	cmd := exec.Command(editor, EnvFilePath)

	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		return errMsg(err)
	})
}
