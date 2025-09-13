package app

import (
	"os"

	c "github.com/Cybernetics354/mayohttp/internal/app/common"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/event"
	"github.com/Cybernetics354/mayohttp/internal/app/val"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *State) Setup() (tea.Model, tea.Cmd) {
	// check whether the config folder exists
	if _, err := os.Stat(val.ConfigFolder); os.IsNotExist(err) {
		err = os.MkdirAll(val.ConfigFolder, 0o755)
		if err != nil {
			return m, c.SendMsg(event.Err(err))
		}
	}

	var saveCmd tea.Cmd
	if _, err := os.Stat(val.DefaultSessionPath); os.IsNotExist(err) {
		saveCmd = c.SendMsg(event.SaveSession{Path: val.DefaultSessionPath})
	}

	return m, tea.Batch(
		saveCmd,
		c.SendMsg(event.CheckEnvFile{}),
		c.SendMsg(event.LoadSession{Path: val.DefaultSessionPath}),
		c.SendMsg(event.RefreshState{}),
		event.ListenResponseCmd(m.resSub),
		event.ListenPipeResponseCmd(m.pipeResSub),
	)
}

func (m *State) HandleErrorMsg(msg event.Err) (tea.Model, tea.Cmd) {
	c.ErrLog.Error(msg.Error())
	return m, nil
}

func (m *State) HandleWindowChange(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	w, h := val.AppStyle.GetFrameSize()
	cw := msg.Width - w
	ch := msg.Height - h

	m.sw = cw
	m.sh = ch

	return m, c.SendMsg(event.RecalculateComponentSizes{})
}

