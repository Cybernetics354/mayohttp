package app

import tea "charm.land/bubbletea/v2"

func (m State) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		sendMsg(setupMsg{}),
	)
}
