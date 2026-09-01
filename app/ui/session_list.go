package ui

import "charm.land/bubbles/v2/list"

func SessionList() list.Model {
	i := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	i.Title = "Session List"
	i.SetShowHelp(false)
	return i
}
