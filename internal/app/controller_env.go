package app

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	c "github.com/Cybernetics354/mayohttp/internal/app/common"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/event"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/intf"
	"github.com/Cybernetics354/mayohttp/internal/app/val"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *State) OpenEnv() (tea.Model, tea.Cmd) {
	editor := c.GetDefaultEditor()
	cmd := exec.Command(editor, val.EnvFilePath)

	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return event.Err(err)
		}

		return nil
	})
}

func (m *State) CheckOrCreateEnvFile() (tea.Model, tea.Cmd) {
	if _, err := os.Stat(val.EnvFilePath); err == nil {
		return m, c.SendMsg(event.Err(err))
	}

	f, err := os.Create(val.EnvFilePath)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return m, tea.Batch(c.SendMsg(event.Err(err)), tea.Quit)
	}

	defer f.Close()
	return m, nil
}

func (m *State) SelectEnv() (tea.Model, tea.Cmd) {
	if m.envList.FilterState() == list.Filtering {
		var cmd tea.Cmd
		m.envList, cmd = m.envList.Update(tea.KeyMsg{
			Type: tea.KeyEnter,
		})
		return m, cmd
	}

	file, ok := m.envList.SelectedItem().(intf.FileItem)
	if !ok {
		return m, c.SendMsg(event.Err(errors.New("no env selected")))
	}

	val.EnvFilePath = file.Name

	return m, c.SendMsg(event.PopStack{})
}

func (m *State) RefreshSelectEnv() (tea.Model, tea.Cmd) {
	files, err := os.ReadDir(".")
	if err != nil {
		return m, c.SendMsg(event.Err(err))
	}

	var fileItems []list.Item
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileItems = append(fileItems, intf.FileItem{Name: file.Name(), Path: file.Name()})
	}

	m.envList.SetItems(fileItems)
	m.envList.KeyMap.Quit.SetEnabled(false)
	return m, nil
}
