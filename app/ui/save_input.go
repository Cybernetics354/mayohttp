package ui

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
)

func SaveInput() textinput.Model {
	i := textinput.New()
	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).PaddingRight(2)
	i.Prompt = "Save as"
	s := i.Styles()
	s.Focused.Prompt = promptStyle
	s.Blurred.Prompt = promptStyle
	i.SetStyles(s)

	return i
}
