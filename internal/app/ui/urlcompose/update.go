package urlcompose

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	m.input, cmd = m.input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m, cmd = m.RunCommand()
		case "ctrl+d":
			m.ClearInput()
		}
	}

	return m, cmd
}
