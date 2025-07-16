package app

import (
	"slices"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *State) HandleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if slices.Contains(homeLayout, m.state) {
		switch {
		case key.Matches(msg, homeMapping.Open):
			return m, sendMsg(openEditorMsg{state: m.state})
		case key.Matches(msg, homeMapping.Commands):
			return m, sendMsg(addStackMsg{state: STATE_COMMAND_PALLETE})
		case key.Matches(msg, homeMapping.Method):
			return m, sendMsg(addStackMsg{state: STATE_METHOD_PALLETE})
		case key.Matches(msg, homeMapping.Next):
			return m, sendMsg(nextSectionMsg{})
		case key.Matches(msg, homeMapping.Back):
			return m, sendMsg(prevSectionMsg{})
		}
	}

	switch {
	case key.Matches(msg, homeMapping.Quit):
		return m.Quit()
	case key.Matches(msg, homeMapping.Run):
		switch m.state {
		case STATE_FOCUS_URL:
			return m, sendMsg(runRequestMsg{})
		case STATE_FOCUS_PIPE:
			return m, sendMsg(runPipeMsg{})
		case STATE_COMMAND_PALLETE:
			return m, sendMsg(selectCommandPalleteMsg{})
		case STATE_METHOD_PALLETE:
			return m, sendMsg(selectMethodPalleteMsg{})
		case STATE_SELECT_SESSION, STATE_SAVE_SESSION:
			return m, sendMsg(selectSessionItemMsg{})
		case STATE_SELECT_ENV:
			return m, sendMsg(selectEnvMsg{})
		}
	}

	switch m.state {
	case STATE_FOCUS_PIPE:
		m.pipe, cmd = m.pipe.Update(msg)
	case STATE_FOCUS_URL:
		m.url, cmd = m.url.Update(msg)
	case STATE_FOCUS_RESPONSE:
		m.response, cmd = m.response.Update(msg)
	case STATE_FOCUS_PIPEDRESP:
		m.pipedresp, cmd = m.pipedresp.Update(msg)
	case STATE_COMMAND_PALLETE:
		m.commands, cmd = m.commands.Update(msg)
	case STATE_METHOD_PALLETE:
		m.methodSelect, cmd = m.methodSelect.Update(msg)
	case STATE_FOCUS_RESPONSE_FILTER:
		m.resFilter, cmd = m.resFilter.HandleKeyPress(msg)
	case STATE_SELECT_ENV:
		m.envList, cmd = m.envList.Update(msg)
	case STATE_SELECT_SESSION, STATE_SAVE_SESSION:
		m.sessionList, cmd = m.sessionList.Update(msg)
	}

	return m, cmd
}
