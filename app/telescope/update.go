package telescope

import (
	tea "charm.land/bubbletea/v2"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+j", "down":
			m.list, cmd = m.list.Update(tea.KeyPressMsg{Code: tea.KeyDown})
			cmds = append(cmds, cmd)
		case "ctrl+k", "up":
			m.list, cmd = m.list.Update(tea.KeyPressMsg{Code: tea.KeyUp})
			cmds = append(cmds, cmd)
		case "ctrl+d":
			m.Clear()
		case "enter":
			cmds = append(cmds, m.GetSelectedMsg())
		}
	}

	m.search, cmd = m.search.Update(msg)
	cmds = append(cmds, cmd)

	m.Sync()

	return m, tea.Batch(cmds...)
}
