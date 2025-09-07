package event

import tea "github.com/charmbracelet/bubbletea"

func ListenResponseCmd(sub chan RequestResult) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func ListenPipeResponseCmd(sub chan PipeResult) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}
