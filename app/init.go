package app

import tea "github.com/charmbracelet/bubbletea"

func (m State) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("MayoHTTP"),
		tea.EnterAltScreen,
		m.spinner.Tick,
		sendMsg(setupMsg{}),
	)
}
