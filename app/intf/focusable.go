package intf

import tea "charm.land/bubbletea/v2"

type IFocusable interface {
	Focus() tea.Cmd
}
