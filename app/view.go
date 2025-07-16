package app

import (
	"fmt"

	"github.com/Cybernetics354/mayohttp/app/ui"
	"github.com/charmbracelet/lipgloss"
)

func (m State) View() string {
	return m.Render()
}

func (m *State) RefreshView() {
	w, h := m.sw, m.sh

	lh := h - 1

	m.help.Width = w
	m.url.Width = w - 5 - len(m.url.Prompt)
	m.pipe.Width = w - 11
	m.response.SetWidth(w)
	m.response.SetHeight(h - 10)
	m.pipedresp.SetWidth(w)
	m.pipedresp.SetHeight(h - 9)
	m.commands.SetSize(ui.ListWidth, lh)
	m.envList.SetSize(ui.ListWidth, lh)
	m.sessionList.SetSize(ui.ListWidth, lh)
	m.methodSelect.SetSize(w, h)
}

func (m *State) Render() string {
	var str string
	switch m.state {
	case STATE_COMMAND_PALLETE:
		str = m.RenderWithListHelp(m.RenderCommandPallete())
	case STATE_METHOD_PALLETE:
		str = lipgloss.JoinVertical(lipgloss.Top, m.methodSelect.View())
	case STATE_SELECT_ENV:
		str = m.RenderWithListHelp(m.RenderEnvList())
	case STATE_SELECT_SESSION, STATE_SAVE_SESSION, STATE_SAVE_SESSION_INPUT:
		str = m.RenderWithListHelp(m.RenderSessionList())
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

func (m *State) PreviewSize() (int, int) {
	return m.sw - ui.ListPreviewWidthMargin, m.sh - 7
}

func (m *State) RenderSessionList() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		ui.ListContainer.Render(m.sessionList.View()),
		" ",
		m.RenderSessionListPreview(),
	)
}

func (m *State) RenderSessionListPreview() string {
	pw, ph := m.PreviewSize()
	item, ok := m.sessionList.SelectedItem().(SessionItem)
	if !ok {
		return ""
	}

	str := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, item.session.Method, " : ", item.session.Url),
		lipgloss.JoinHorizontal(lipgloss.Top, "PIPE : ", item.session.Pipe),
		item.session.PipedResponse,
	)

	preview := ui.Preview{
		Header:    "Preview",
		Width:     pw,
		MaxHeight: ph,
		Body:      str,
	}

	return preview.Render()
}

func (m *State) RenderEnvList() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		ui.ListContainer.Render(m.envList.View()),
		" ",
		m.RenderEnvListPreview(),
	)
}

func (m *State) RenderEnvListPreview() string {
	pw, ph := m.PreviewSize()
	file, ok := m.envList.SelectedItem().(fileItem)

	prev := ""
	if ok {
		uiPrev := ui.Preview{
			Header:    "Preview",
			Width:     pw,
			MaxHeight: ph,
			Body:      printval(file.path, true),
		}

		prev = lipgloss.JoinVertical(
			lipgloss.Left,
			uiPrev.Render(),
		)
	}

	return prev
}

func (m *State) RenderCommandPallete() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		ui.ListContainer.Render(m.commands.View()),
		" ",
		m.RenderCommandPalletePreview(),
	)
}

func (m *State) RenderCommandPalletePreview() string {
	pw, ph := m.PreviewSize()
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

	preview := ui.Preview{
		Header:    "Preview",
		Width:     pw,
		MaxHeight: ph,
		Body:      str,
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		preview.Render(),
	)
}

func (m *State) RenderHelp() string {
	return m.help.View(m.keys)
}

func (m *State) RenderWithListHelp(body string) string {
	return lipgloss.JoinVertical(lipgloss.Left, body, m.help.View(listMapping))
}

func (m *State) RenderURL() string {
	c := m.url.View()

	if m.state == STATE_FOCUS_URL {
		return ui.FocusInputContainer.Render(c)
	}

	return ui.BlurInputContainer.Render(c)
}

func (m *State) RenderPipe() string {
	c := m.pipe.View()

	if m.state == STATE_FOCUS_PIPE {
		return ui.FocusInputContainer.Render(c)
	}

	return ui.BlurInputContainer.Render(c)
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
