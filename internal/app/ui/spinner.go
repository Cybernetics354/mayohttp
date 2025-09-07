package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

func Spinner() spinner.Model {
	i := spinner.New()
	i.Spinner = spinner.Dot
	i.Style = lipgloss.NewStyle().Foreground(FocusColor)
	return i
}
