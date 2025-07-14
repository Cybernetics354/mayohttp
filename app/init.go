package app

import tea "github.com/charmbracelet/bubbletea"

func (m State) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("MayoHTTP"),
		tea.EnterAltScreen,
		m.spinner.Tick,
		sendMsg(checkEnvFileMsg{}),
		sendMsg(loadSessionMsg{path: defaultSessionPath}),
		sendMsg(refreshStateMsg{}),
		listenResponseCmd(m.resSub),
		listenPipeResponseCmd(m.pipeResSub),
	)
}
