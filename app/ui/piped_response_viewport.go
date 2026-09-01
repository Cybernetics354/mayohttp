package ui

import (
	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
)

func PipedResponseViewport() viewport.Model {
	i := viewport.New()
	i.Style.BorderStyle(lipgloss.RoundedBorder())
	i.Style.BorderForeground(FocusColor)
	return i
}
