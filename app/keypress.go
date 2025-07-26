package app

import (
	"errors"
	"slices"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *State) HandleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if slices.Contains(homeLayout, m.state) {
		switch {
		case key.Matches(msg, homeMapping.QuickAccess):
			return m, sendMsg(openTelescopeMsg{teleType: TELESCOPE_QUICK_ACCESS})
		case key.Matches(msg, homeMapping.ComposeUrl):
			m.urlcompose.SetUrl(m.url.Value())
			return m, sendMsg(addStackMsg{state: STATE_URL_COMPOSE})
		case key.Matches(msg, homeMapping.CopyToClipboard):
			return m, sendMsg(copyToClipboardMsg{})
		case key.Matches(msg, homeMapping.OpenEnv):
			return m, sendMsg(openEnvMsg{})
		case key.Matches(msg, homeMapping.Open):
			return m, sendMsg(openEditorMsg{state: m.state})
		case key.Matches(msg, homeMapping.Keybinding):
			return m, sendMsg(addStackMsg{state: STATE_KEYBINDING_MODAL})
		case key.Matches(msg, homeMapping.Commands):
			return m, sendMsg(addStackMsg{state: STATE_COMMAND_PALLETE})
		case key.Matches(msg, homeMapping.Method):
			return m, sendMsg(openTelescopeMsg{teleType: TELESCOPE_METHOD_PALLETE})
		case key.Matches(msg, homeMapping.Next):
			return m, sendMsg(nextSectionMsg{})
		case key.Matches(msg, homeMapping.Back):
			return m, sendMsg(prevSectionMsg{})
		case key.Matches(msg, homeMapping.Save):
			return m, tea.Sequence(
				sendMsg(addStackMsg{state: STATE_SAVE_SESSION}),
				sendMsg(loadSessionListMsg{}),
			)
		case key.Matches(msg, homeMapping.OpenSession):
			return m, tea.Sequence(
				sendMsg(addStackMsg{state: STATE_SELECT_SESSION}),
				sendMsg(loadSessionListMsg{}),
			)
		}
	}

	switch {
	case key.Matches(msg, homeMapping.Keybinding):
		// if the keybinding modal is open, then close it
		if m.state == STATE_KEYBINDING_MODAL {
			return m, sendMsg(popStackMsg{})
		}
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
		case STATE_SAVE_SESSION_INPUT, STATE_SESSION_RENAME_INPUT:
			return m, sendMsg(saveInputSubmitMsg{})
		}
	}

	switch m.state {
	case STATE_FOCUS_PIPE:
		m.pipe, cmd = m.pipe.Update(msg)
	case STATE_TELESCOPE:
		m.telescope, cmd = m.telescope.Update(msg)
	case STATE_URL_COMPOSE:
		m.urlcompose, cmd = m.urlcompose.Update(msg)
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
	case STATE_SESSION_RENAME_INPUT, STATE_SAVE_SESSION_INPUT:
		m.saveInput, cmd = m.saveInput.Update(msg)
	case STATE_SAVE_SESSION:
		m.sessionList, cmd = m.sessionList.Update(msg)
		if m.sessionList.FilterState() != list.Filtering {
			switch {
			case key.Matches(msg, saveListMapping.New):
				m.saveInput.SetValue("")
				m.saveInput.Prompt = "New Session Name"
				return m, sendMsg(addStackMsg{state: STATE_SAVE_SESSION_INPUT})
			}
		}
	case STATE_SELECT_SESSION:
		m.sessionList, cmd = m.sessionList.Update(msg)
		if m.sessionList.FilterState() != list.Filtering {
			switch {
			case key.Matches(msg, sessionListMapping.Rename):
				i, ok := m.sessionList.SelectedItem().(SessionItem)
				if !ok {
					return m, sendMsg(errMsg(errors.New("no session selected")))
				}

				m.saveInput.Prompt = "Rename Session"
				m.saveInput.SetValue(i.Title())

				return m, sendMsg(addStackMsg{state: STATE_SESSION_RENAME_INPUT})
			case key.Matches(msg, sessionListMapping.Delete):
				return m, sendMsg(deleteSessionItemMsg{})
			}
		}
	}

	return m, cmd
}
