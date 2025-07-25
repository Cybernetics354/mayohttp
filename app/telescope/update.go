package telescope

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+j", "down":
			m.list, cmd = m.list.Update(tea.KeyMsg{Type: tea.KeyDown})
			cmds = append(cmds, cmd)
		case "ctrl+k", "up":
			m.list, _ = m.list.Update(tea.KeyMsg{Type: tea.KeyUp})
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
