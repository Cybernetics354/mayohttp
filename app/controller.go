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
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type requestResponse struct {
	err error
	res string
}

type requestPipeResponse struct {
	err error
	res string
}

func listenResponse(sub chan requestResponse) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func listenPipeResponse(sub chan requestPipeResponse) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func (m State) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("Mayo HTTP"),
		tea.EnterAltScreen,
		m.spinner.Tick,
		listenResponse(m.resSub),
		listenPipeResponse(m.pipeResSub),
	)
}

func (m State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.HandleWindowChange(msg)
	case tea.KeyMsg:
		return m.HandleKeyPress(msg)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case requestResponse:
		m.HideSpinner()
		m.response.SetValue(msg.res)
		m.RunPipe()
		return m, listenResponse(m.resSub)
	case requestPipeResponse:
		m.HideSpinner()
		m.pipedresp.SetValue(msg.res)
		return m, listenPipeResponse(m.pipeResSub)
	case errMsg:
		return m.HandleErrorMsg(msg)
	}

	// for filtering, i don't know why it doesn't count as tea.KeyMsg
	// i don't really know the exact type too so i'll just put it on default handler
	//
	// (i've wasted my 2 hours for this shit)
	var cmd tea.Cmd
	m.commands, cmd = m.commands.Update(msg)
	return m, cmd
}

func (m State) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	return m.Render()
}

func (m *State) Request() {
	m.ShowSpinner()
	go func() {
		req, err := http.NewRequest(m.method, m.url.Value(), nil)
		if err != nil {
			m.resSub <- requestResponse{
				err: err,
				res: fmt.Sprintf("Request error : %s", err.Error()),
			}
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			m.resSub <- requestResponse{
				err: err,
				res: fmt.Sprintf("Request do error : %s", err.Error()),
			}
			return
		}

		defer resp.Body.Close()

		respDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			m.resSub <- requestResponse{
				err: err,
				res: fmt.Sprintf("Response dump error : %s", err.Error()),
			}
			return
		}

		m.resSub <- requestResponse{res: string(respDump)}
	}()
}

func (m *State) RunPipe() {
	m.ShowSpinner()
	go func() {
		resp, pipe := m.response.Value(), m.pipe.Value()
		if resp == "" || pipe == "" {
			m.pipeResSub <- requestPipeResponse{res: m.response.Value()}
			return
		}

		command := exec.Command("bash", "-c", fmt.Sprintf("echo '%s' | %s", resp, pipe))
		output, err := command.CombinedOutput()
		if err != nil {
			m.pipeResSub <- requestPipeResponse{
				err: err,
				res: fmt.Sprintf("Command pipe error : %s", err.Error()),
			}
			return
		}

		m.pipeResSub <- requestPipeResponse{res: string(output)}
	}()
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

	m.RecalculateComponentSize()

	return m, nil
}

func (m *State) HandleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, keyMaps.Open):
		return m, m.OpenOnEditor()
	case key.Matches(msg, keyMaps.Quit):
		m.Quit()
		if m.quitting {
			return m, tea.Quit
		}
		return m, nil
	case key.Matches(msg, keyMaps.Commands):
		m.SetState(COMMAND_PALLETE)
		return m, nil
	case key.Matches(msg, keyMaps.Next):
		m.NextState()
		return m, nil
	case key.Matches(msg, keyMaps.Back):
		m.PrevState()
		return m, nil
	case key.Matches(msg, keyMaps.Run):
		if m.state == FOCUS_URL {
			m.Request()
			return m, nil
		}

		if m.state == FOCUS_PIPE {
			m.RunPipe()
			return m, nil
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

func (m *State) NextState() {
	if !slices.Contains(homeLayout, m.state) {
		return
	}

	var nextState string
	for i, state := range homeLayout {
		if state != m.state {
			continue
		}

		nextState = homeLayout[min(len(homeLayout)-1, i+1)]
		break
	}

	m.SetState(nextState)
}

func (m *State) PrevState() {
	if !slices.Contains(homeLayout, m.state) {
		return
	}

	var prevState string
	for i, state := range homeLayout {
		if state != m.state {
			continue
		}

		prevState = homeLayout[max(0, i-1)]
		break
	}

	m.SetState(prevState)
}

func (m *State) Quit() {
	if m.state == FOCUS_URL {
		m.quitting = true
		return
	}

	m.SetState(FOCUS_URL)
}

func (m *State) SetState(state string) {
	m.state = state
	m.RefreshState()
}

func (m *State) RefreshState() {
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
}

func (m *State) OpenOnEditor() tea.Cmd {
	var str string
	switch f := m.GetFocusedField().(type) {
	case *textarea.Model:
		str = f.Value()
	case *textinput.Model:
		str = f.Value()
	default:
		return nil
	}

	dir := filepath.Dir(tempFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			m.err = err
			return nil
		}
	}

	f, err := os.Create(tempFilePath)
	if err != nil {
		m.err = err
		return nil
	}

	f.WriteString(str)
	f.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	cmd := exec.Command(editor, tempFilePath)
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
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

func (m *State) ShowSpinner() {
	m.showSpinner = true
}

func (m *State) HideSpinner() {
	m.showSpinner = false
}
