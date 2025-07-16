package ui

import "github.com/charmbracelet/lipgloss"

var (
	BorderStyle = lipgloss.RoundedBorder()
	FocusColor  = lipgloss.Color("205")
	BlurColor   = lipgloss.Color("243")

	ListWidth              = 35
	ListPreviewWidthMargin = ListWidth + 5
	ListContainer          = lipgloss.NewStyle().Width(ListWidth)

	FocusInputContainer = lipgloss.NewStyle().
				Border(BorderStyle).
				BorderForeground(FocusColor).Padding(0, 1)
	BlurInputContainer = lipgloss.NewStyle().
				Border(BorderStyle).
				BorderForeground(BlurColor).Padding(0, 1)
	FocusTextarea = lipgloss.NewStyle().
			Border(BorderStyle).
			BorderForeground(FocusColor)
	BlurTextarea = lipgloss.NewStyle().
			Border(BorderStyle).
			BorderForeground(BlurColor)
	FocusTextareaLineNumber = lipgloss.NewStyle().Foreground(BlurColor)
	BlurTextareaLineNumber  = lipgloss.NewStyle().Foreground(BlurColor)

	UrlPromptStyle = lipgloss.NewStyle().Foreground(FocusColor)

	PreviewHeaderStyle = lipgloss.NewStyle().
				Background(FocusColor).
				Padding(0, 1).
				MarginLeft(1)
	PreviewBodyStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(FocusColor).
				Padding(0, 1)
)
