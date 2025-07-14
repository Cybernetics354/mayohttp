package ui

import "github.com/charmbracelet/bubbles/list"

func CommandList(items []list.Item) list.Model {
	i := list.New(items, list.NewDefaultDelegate(), 0, 0)
	i.Title = "Commands Pallete"
	i.KeyMap.Quit.SetEnabled(false)
	return i
}
