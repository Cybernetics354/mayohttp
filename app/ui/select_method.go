package ui

import "charm.land/bubbles/v2/list"

func SelectMethod(items []list.Item) list.Model {
	i := list.New(items, list.NewDefaultDelegate(), 0, 0)
	i.Title = "Select Method"
	i.SetShowHelp(false)
	return i
}
