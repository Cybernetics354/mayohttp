package app

import (
	"fmt"
	"net/http"
	"os/exec"
	"slices"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m State) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.onWindowChanged(msg)
	case tea.KeyMsg:
		return m.onKeyPressed(msg)
	case errMsg:
		return m.onRenderError(msg)
	}

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
	req, err := http.NewRequest(m.method, m.url.Value(), nil)
	if err != nil {
		m.response.SetValue(err.Error())
		return
	}

	m.response.SetValue(formatResponse(req))
	m.RunPipe()
}

func (m *State) RunPipe() {
	resp := m.response.Value()
	pipe := m.pipe.Value()
	if resp == "" || pipe == "" {
		m.pipedresp.SetValue(m.response.Value())
		return
	}

	command := exec.Command("bash", "-c", fmt.Sprintf("echo '%s' | %s", resp, pipe))
	output, err := command.Output()
	if err != nil {
		m.pipedresp.SetValue(err.Error())
		return
	}

	m.pipedresp.SetValue(string(output))
}

func (m *State) onRenderError(msg errMsg) (tea.Model, tea.Cmd) {
	m.err = msg
	return m, nil
}

func (m *State) onWindowChanged(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	w, h := appStyle.GetFrameSize()
	cw := msg.Width - w
	ch := msg.Height - h

	m.sw = cw
	m.sh = ch

	m.RecalculateComponentSize()

	return m, nil
}

func (m *State) onKeyPressed(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, mapped.Quit):
		m.handleQuit()
		if m.quitting {
			return m, tea.Quit
		}
		return m, nil
	case key.Matches(msg, mapped.Commands):
		m.state = COMMAND_PALLETE
		m.setupState()
		return m, nil
	case key.Matches(msg, mapped.Next):
		m.nextState()
		return m, nil
	case key.Matches(msg, mapped.Back):
		m.prevState()
		return m, nil
	case key.Matches(msg, mapped.Run):
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
	case COMMAND_PALLETE:
		m.commands, cmd = m.commands.Update(msg)
	}

	return m, cmd
}

func (m *State) nextState() {
	if !slices.Contains(homeLayout, m.state) {
		return
	}

	var nextState string
	for i, state := range homeLayout {
		if state != m.state {
			continue
		}
		if i == len(homeLayout)-1 {
			nextState = homeLayout[0]
			break
		}
		nextState = homeLayout[i+1]
		break
	}
	m.state = nextState
	m.setupState()
}

func (m *State) prevState() {
	if !slices.Contains(homeLayout, m.state) {
		return
	}

	var prevState string
	for i, state := range homeLayout {
		if state != m.state {
			continue
		}
		if i == 0 {
			prevState = homeLayout[len(homeLayout)-1]
			break
		}
		prevState = homeLayout[i-1]
		break
	}
	m.state = prevState
	m.setupState()
}

func (m *State) handleQuit() {
	if m.state == FOCUS_URL {
		m.quitting = true
		return
	}

	m.state = FOCUS_URL
	m.setupState()
}

func (m *State) setupState() {
	m.url.Blur()
	m.response.Blur()
	m.pipe.Blur()
	m.pipedresp.Blur()
	m.body.Blur()
	m.header.Blur()
	m.commands.ResetFilter()

	switch m.state {
	case FOCUS_URL:
		m.url.Focus()
	case FOCUS_RESPONSE:
		m.response.Focus()
	case FOCUS_PIPE:
		m.pipe.Focus()
	case FOCUS_PIPEDRESP:
		m.pipedresp.Focus()
	case EDIT_BODY:
		m.body.Focus()
	case EDIT_HEADER:
		m.header.Focus()
	case COMMAND_PALLETE:
		m.commands.FilterInput.Focus()
	}
}
