package app

import (
	"errors"
	"slices"

	c "github.com/Cybernetics354/mayohttp/internal/app/common"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/event"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/store"
	"github.com/Cybernetics354/mayohttp/internal/app/val"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *State) HandleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if slices.Contains(val.HomeLayout, m.state) {
		switch {
		case key.Matches(msg, val.HomeMapping.QuickAccess):
			return m, c.SendMsg(event.OpenTelescope{TeleType: val.TELESCOPE_QUICK_ACCESS})
		case key.Matches(msg, val.HomeMapping.ComposeUrl):
			m.urlcompose.SetUrl(m.url.Value())
			return m, c.SendMsg(event.AddStack{State: val.STATE_URL_COMPOSE})
		case key.Matches(msg, val.HomeMapping.CopyToClipboard):
			return m, c.SendMsg(event.CopyToClipboard{})
		case key.Matches(msg, val.HomeMapping.ClearInput):
			return m.ClearFocusedInput()
		case key.Matches(msg, val.HomeMapping.OpenEnv):
			return m, c.SendMsg(event.OpenEnv{})
		case key.Matches(msg, val.HomeMapping.Open):
			return m, c.SendMsg(event.OpenEditor{State: m.state})
		case key.Matches(msg, val.HomeMapping.Keybinding):
			return m, c.SendMsg(event.AddStack{State: val.STATE_KEYBINDING_MODAL})
		case key.Matches(msg, val.HomeMapping.Commands):
			return m, c.SendMsg(event.AddStack{State: val.STATE_COMMAND_PALLETE})
		case key.Matches(msg, val.HomeMapping.Method):
			return m, c.SendMsg(event.OpenTelescope{TeleType: val.TELESCOPE_METHOD_PALLETE})
		case key.Matches(msg, val.HomeMapping.Next):
			return m, c.SendMsg(event.NextSection{})
		case key.Matches(msg, val.HomeMapping.Back):
			return m, c.SendMsg(event.PrevSection{})
		case key.Matches(msg, val.HomeMapping.Save):
			return m, tea.Sequence(
				c.SendMsg(event.AddStack{State: val.STATE_SAVE_SESSION}),
				c.SendMsg(event.LoadSessionList{}),
			)
		case key.Matches(msg, val.HomeMapping.OpenSession):
			return m, tea.Sequence(
				c.SendMsg(event.LoadSessionList{}),
				c.SendMsg(event.AddStack{State: val.STATE_SELECT_SESSION}),
			)
		}
	}

	switch {
	case key.Matches(msg, val.HomeMapping.Keybinding):
		// if the keybinding modal is open, then close it
		if m.state == val.STATE_KEYBINDING_MODAL {
			return m, c.SendMsg(event.PopStack{})
		}
	case key.Matches(msg, val.HomeMapping.Quit):
		return m.Quit()
	case key.Matches(msg, val.HomeMapping.Run):
		switch m.state {
		case val.STATE_FOCUS_URL:
			return m, c.SendMsg(event.RunRequest{})
		case val.STATE_FOCUS_PIPE:
			return m, c.SendMsg(event.RunPipe{})
		case val.STATE_COMMAND_PALLETE:
			return m, c.SendMsg(event.SelectCommandPallete{})
		case val.STATE_METHOD_PALLETE:
			return m, c.SendMsg(event.SelectMethodPallete{})
		case val.STATE_SELECT_SESSION, val.STATE_SAVE_SESSION:
			return m, c.SendMsg(event.SelectSessionItem{})
		case val.STATE_SELECT_ENV:
			return m, c.SendMsg(event.SelectEnv{})
		case val.STATE_SAVE_SESSION_INPUT, val.STATE_SESSION_RENAME_INPUT:
			return m, c.SendMsg(event.SaveInputSubmit{})
		}
	}

	switch m.state {
	case val.STATE_FOCUS_PIPE:
		m.pipe, cmd = m.pipe.Update(msg)
	case val.STATE_TELESCOPE:
		m.telescope, cmd = m.telescope.Update(msg)
	case val.STATE_URL_COMPOSE:
		m.urlcompose, cmd = m.urlcompose.Update(msg)
	case val.STATE_FOCUS_URL:
		m.url, cmd = m.url.Update(msg)
	case val.STATE_FOCUS_RESPONSE:
		m.response, cmd = m.response.Update(msg)
	case val.STATE_FOCUS_PIPEDRESP:
		m.pipedresp, cmd = m.pipedresp.Update(msg)
	case val.STATE_COMMAND_PALLETE:
		m.commands, cmd = m.commands.Update(msg)
	case val.STATE_METHOD_PALLETE:
		m.methodSelect, cmd = m.methodSelect.Update(msg)
	case val.STATE_FOCUS_RESPONSE_FILTER:
		m.resFilter, cmd = m.resFilter.HandleKeyPress(msg)
	case val.STATE_SELECT_ENV:
		m.envList, cmd = m.envList.Update(msg)
	case val.STATE_SESSION_RENAME_INPUT, val.STATE_SAVE_SESSION_INPUT:
		m.saveInput, cmd = m.saveInput.Update(msg)
	case val.STATE_SAVE_SESSION:
		m.sessionList, cmd = m.sessionList.Update(msg)
		if m.sessionList.FilterState() == list.Filtering {
			return m, cmd
		}

		switch {
		case key.Matches(msg, val.SaveListMapping.New):
			m.saveInput.SetValue("")
			m.saveInput.Prompt = "New Session Name"
			return m, c.SendMsg(event.AddStack{State: val.STATE_SAVE_SESSION_INPUT})
		}
	case val.STATE_SELECT_SESSION:
		m.sessionList, cmd = m.sessionList.Update(msg)
		if m.sessionList.FilterState() == list.Filtering {
			return m, cmd
		}

		switch {
		case key.Matches(msg, val.SessionListMapping.Delete):
			return m, c.SendMsg(event.DeleteSessionItem{})
		case key.Matches(msg, val.SessionListMapping.Rename):
			i, ok := m.sessionList.SelectedItem().(*store.Session)
			if !ok {
				return m, c.SendMsg(event.Err(errors.New("no session selected")))
			}

			m.saveInput.Prompt = "Rename Session"
			m.saveInput.SetValue(i.Title())

			return m, c.SendMsg(event.AddStack{State: val.STATE_SESSION_RENAME_INPUT})
		}
	}

	return m, cmd
}
