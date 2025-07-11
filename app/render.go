package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *State) RecalculateComponentSize() (tea.Model, tea.Cmd) {
	w, h := m.sw, m.sh

	m.help.Width = w
	m.url.Width = w - 5 - len(m.url.Prompt)
	m.pipe.Width = w - 11
	m.response.SetWidth(w)
	m.response.SetHeight(h - 8)
	m.pipedresp.SetWidth(w)
	m.pipedresp.SetHeight(h - 8)
	m.commands.SetSize(w, h)
	m.methodSelect.SetSize(w, h)

	return m, nil
}

func (m *State) Render() string {
	var str string
	switch m.state {
	case COMMAND_PALLETE:
		str = m.commands.View()
	case METHOD_PALLETE:
		str = m.methodSelect.View()
	default:
		str = lipgloss.JoinVertical(
			lipgloss.Top,
			m.RenderURL(),
			m.RenderPipe(),
			m.RenderPipedResponse(),
			lipgloss.JoinHorizontal(lipgloss.Center, m.RenderSpinner(), m.RenderHelp()),
		)
		break
	}

	return appStyle.Render(str)
}

func (m *State) RenderHelp() string {
	return m.help.View(m.keys)
}

func (m *State) RenderURL() string {
	c := m.url.View()

	if m.state == FOCUS_URL {
		return focusInputContainer.Render(c)
	}

	return blurInputContainer.Render(c)
}

func (m *State) RenderPipe() string {
	c := m.pipe.View()

	if m.state == FOCUS_PIPE {
		return focusInputContainer.Render(c)
	}

	return blurInputContainer.Render(c)
}

func (m *State) RenderResponse() string {
	return m.response.View()
}

func (m *State) RenderPipedResponse() string {
	return m.pipedresp.View()
}

func (m *State) RenderSpinner() string {
	if !m.showSpinner {
		return ""
	}

	return fmt.Sprintf("%s ", m.spinner.View())
}
