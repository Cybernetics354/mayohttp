package ui

import (
	"charm.land/bubbles/v2/spinner"
	"charm.land/lipgloss/v2"
)

func Spinner() spinner.Model {
	i := spinner.New()
	i.Spinner = spinner.Dot
	i.Style = lipgloss.NewStyle().Foreground(FocusColor)
	return i
}
