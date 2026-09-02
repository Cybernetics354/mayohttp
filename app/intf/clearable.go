package intf

import tea "charm.land/bubbletea/v2"

type IClearable interface {
	Clear() tea.Cmd
}
