package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func PipeInput() textinput.Model {
	i := textinput.New()
	i.PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).PaddingRight(2)
	i.Prompt = "PIPE"

	return i
}
