package app

import (
	c "github.com/Cybernetics354/mayohttp/internal/app/common"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/event"
	tea "github.com/charmbracelet/bubbletea"
)

func (m State) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("MayoHTTP"),
		tea.EnterAltScreen,
		m.spinner.Tick,
		c.SendMsg(event.Setup{}),
	)
}
