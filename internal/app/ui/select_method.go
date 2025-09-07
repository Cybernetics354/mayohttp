package ui

import "github.com/charmbracelet/bubbles/list"

func SelectMethod(items []list.Item) list.Model {
	i := list.New(items, list.NewDefaultDelegate(), 0, 0)
	i.Title = "Select Method"
	i.SetShowHelp(false)
	return i
}
