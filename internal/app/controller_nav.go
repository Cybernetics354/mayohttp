package app

import (
	"slices"

	c "github.com/Cybernetics354/mayohttp/internal/app/common"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/event"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/intf"
	"github.com/Cybernetics354/mayohttp/internal/app/val"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *State) NextSection() (tea.Model, tea.Cmd) {
	if !slices.Contains(val.HomeLayout, m.state) {
		return m, nil
	}

	var nextState string
	for i, state := range val.HomeLayout {
		if state != m.state {
			continue
		}

		nextState = val.HomeLayout[min(len(val.HomeLayout)-1, i+1)]
		break
	}

	return m, c.SendMsg(event.SetState{State: nextState})
}

func (m *State) PrevSection() (tea.Model, tea.Cmd) {
	if !slices.Contains(val.HomeLayout, m.state) {
		return m, nil
	}

	var prevState string
	for i, state := range val.HomeLayout {
		if state != m.state {
			continue
		}

		prevState = val.HomeLayout[max(0, i-1)]
		break
	}

	return m, c.SendMsg(event.SetState{State: prevState})
}

func (m *State) Quit() (tea.Model, tea.Cmd) {
	if m.state == val.STATE_COMMAND_PALLETE && m.commands.FilterState() != list.Unfiltered {
		m.commands.ResetFilter()
		return m, nil
	}

	if m.state == val.STATE_SELECT_ENV && m.envList.FilterState() != list.Unfiltered {
		m.envList.ResetFilter()
		return m, nil
	}

	if m.state == val.STATE_METHOD_PALLETE && m.methodSelect.FilterState() != list.Unfiltered {
		m.methodSelect.ResetFilter()
		return m, nil
	}

	if slices.Contains([]string{val.STATE_SELECT_SESSION, val.STATE_SAVE_SESSION}, m.state) &&
		m.sessionList.FilterState() != list.Unfiltered {
		m.sessionList.ResetFilter()
		return m, nil
	}

	if len(m.stateStack) <= 1 {
		go m.SaveSession(event.SaveSession{Path: val.DefaultSessionPath})
		return m, tea.Quit
	}

	return m, c.SendMsg(event.PopStack{})
}

func (m *State) AddStack(state string) (tea.Model, tea.Cmd) {
	m.stateStack = append(m.stateStack, state)
	m.state = state

	return m, c.SendMsg(event.RefreshState{})
}

func (m *State) PopStack() (tea.Model, tea.Cmd) {
	if len(m.stateStack) <= 1 {
		return m, nil
	}

	m.stateStack = m.stateStack[:len(m.stateStack)-1]
	m.state = m.stateStack[len(m.stateStack)-1]

	return m, c.SendMsg(event.RefreshState{})
}

func (m *State) PopStackRoot() (tea.Model, tea.Cmd) {
	if len(m.stateStack) <= 1 {
		return m, nil
	}

	m.stateStack = []string{m.stateStack[0]}
	m.state = m.stateStack[len(m.stateStack)-1]

	return m, c.SendMsg(event.RefreshState{})
}

func (m *State) SetState(state string) (tea.Model, tea.Cmd) {
	m.state = state
	m.stateStack[len(m.stateStack)-1] = state

	return m, c.SendMsg(event.RefreshState{})
}

func (m *State) RefreshState() (tea.Model, tea.Cmd) {
	m.url.Blur()
	m.response.Blur()
	m.pipe.Blur()
	m.pipedresp.Blur()
	m.body.Blur()
	m.header.Blur()
	m.resFilter.Blur()
	m.saveInput.Blur()

	switch f := m.GetFocusedField().(type) {
	case *textarea.Model:
		f.Focus()
	case *textinput.Model:
		f.Focus()
	case *intf.ResponseFilter:
		f.Focus()
	}

	return m, nil
}
