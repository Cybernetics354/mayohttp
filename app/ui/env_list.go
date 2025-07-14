package ui

import "github.com/charmbracelet/bubbles/list"

func EnvList() list.Model {
	i := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	i.Title = "Select ENV"
	i.KeyMap.Quit.SetEnabled(false)
	return i
}
