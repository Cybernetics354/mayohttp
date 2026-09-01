package ui

import "charm.land/bubbles/v2/list"

func CommandList(items []list.Item) list.Model {
	i := list.New(items, list.NewDefaultDelegate(), 0, 0)
	i.Title = "Commands Pallete"
	i.SetShowHelp(false)
	return i
}
