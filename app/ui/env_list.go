package ui

import "charm.land/bubbles/v2/list"

func EnvList() list.Model {
	i := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	i.Title = "Select ENV"
	i.SetShowHelp(false)
	return i
}
