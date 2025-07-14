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
	m.response.SetHeight(h - 10)
	m.pipedresp.SetWidth(w)
	m.pipedresp.SetHeight(h - 9)
	m.commands.SetSize(30, h)
	m.envList.SetSize(30, h)
	m.methodSelect.SetSize(w, h)

	return m, nil
}

func (m *State) Render() string {
	var str string
	switch m.state {
	case STATE_COMMAND_PALLETE:
		str = m.RenderCommandPallete()
	case STATE_METHOD_PALLETE:
		str = lipgloss.JoinVertical(lipgloss.Top, m.methodSelect.View())
	case STATE_SELECT_ENV:
		str = m.RenderEnvList()
	default:
		str = lipgloss.JoinVertical(
			lipgloss.Top,
			m.RenderURL(),
			m.RenderPipe(),
			lipgloss.NewStyle().PaddingLeft(1).Render(
				m.resFilter.Render(),
			),
			m.RenderPipedResponse(),
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.RenderHelp(),
				"  | ",
				m.RenderSpinner(),
				m.activity,
				" | ",
				fmt.Sprintf("ENV(%s)", EnvFilePath),
			),
		)
	}

	return appStyle.Render(str)
}

func (m *State) RenderEnvList() string {
	prevWidth := m.sw - 45
	file, ok := m.envList.SelectedItem().(fileItem)

	prev := ""
	if ok {
		prev = lipgloss.JoinVertical(
			lipgloss.Left,
			previewHeaderStyle.Render("Preview"),
			previewBodyStyle.Width(prevWidth).
				Render(lipgloss.NewStyle().MaxHeight(m.sh-7).Render(printval(file.path, true))),
		)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.envList.View(),
		" ",
		prev,
	)
}

func (m *State) RenderCommandPallete() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.commands.View(),
		" ",
		m.RenderCommandPalletePreview(),
	)
}

func (m *State) RenderCommandPalletePreview() string {
	prevWidth := m.sw - 45
	var str string

	command, ok := m.commands.SelectedItem().(commandPallete)
	if !ok {
		return ""
	}

	switch command.commandId {
	case COMMAND_OPEN_ENV:
		str = printval(EnvFilePath, true)
	case COMMAND_OPEN_BODY:
		str = printval(m.body.Value(), false)
	case COMMAND_OPEN_HEADER:
		str = printval(m.header.Value(), false)
	case COMMAND_SELECT_METHOD:
		str = fmt.Sprintf("Current method : %s", m.method)
	case COMMAND_CHANGE_ENV:
		str = fmt.Sprintf("Current ENV : %s", EnvFilePath)
	}

	if len(str) <= 0 {
		return ""
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		previewHeaderStyle.
			Render("Preview"),
		previewBodyStyle.
			Width(prevWidth).
			Render(lipgloss.NewStyle().MaxHeight(m.sh-7).Render(str)),
	)
}

func (m *State) RenderHelp() string {
	return m.help.View(m.keys)
}

func (m *State) RenderURL() string {
	c := m.url.View()

	if m.state == STATE_FOCUS_URL {
		return focusInputContainer.Render(c)
	}

	return blurInputContainer.Render(c)
}

func (m *State) RenderPipe() string {
	c := m.pipe.View()

	if m.state == STATE_FOCUS_PIPE {
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
