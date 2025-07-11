package app

import "github.com/charmbracelet/lipgloss"

var (
	borderStyle         = lipgloss.RoundedBorder()
	focusBorderColor    = lipgloss.Color("205")
	blurBorderColor     = lipgloss.Color("243")
	focusInputContainer = lipgloss.NewStyle().
				Border(borderStyle).
				BorderForeground(focusBorderColor).Padding(0, 1)
	blurInputContainer = lipgloss.NewStyle().
				Border(borderStyle).
				BorderForeground(blurBorderColor).Padding(0, 1)
	focusTextarea = lipgloss.NewStyle().
			Border(borderStyle).
			BorderForeground(focusBorderColor)
	blurTextarea = lipgloss.NewStyle().
			Border(borderStyle).
			BorderForeground(blurBorderColor)
	focusTextareaLineNumber = lipgloss.NewStyle().Foreground(blurBorderColor)
	blurTextareaLineNumber  = lipgloss.NewStyle().Foreground(blurBorderColor)

	urlPromptStyle = lipgloss.NewStyle().Foreground(focusBorderColor)
)
