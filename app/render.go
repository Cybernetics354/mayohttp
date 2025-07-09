package app

import "fmt"

func (m *State) RecalculateComponentSize() {
	w, h := m.sw, m.sh

	m.help.Width = w
	m.url.Width = w - 10
	m.pipe.Width = w - 11
	m.response.SetWidth(w)
	m.response.SetHeight(h - 8)
	m.pipedresp.SetWidth(w)
	m.pipedresp.SetHeight(h - 8)
	m.commands.SetSize(w, h)
}

func (m *State) Render() string {
	var str string
	switch m.state {
	case FOCUS_URL, FOCUS_PIPE, FOCUS_PIPEDRESP:
		str = fmt.Sprintf(
			"%s\n%s\n%s\n%s",
			m.RenderURL(),
			m.RenderPipe(),
			m.RenderPipedResponse(),
			m.RenderHelp(),
		)
		break
	case COMMAND_PALLETE:
		str = fmt.Sprintf(
			"%s",
			m.commands.View(),
		)
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
