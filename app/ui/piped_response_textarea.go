package ui

import "github.com/charmbracelet/bubbles/textarea"

func PipedResponseTextarea() textarea.Model {
	i := textarea.New()
	i.ShowLineNumbers = true
	i.Prompt = ""
	i.FocusedStyle.Base = FocusTextarea
	i.BlurredStyle.Base = BlurTextarea
	i.FocusedStyle.LineNumber = FocusTextareaLineNumber
	i.BlurredStyle.LineNumber = BlurTextareaLineNumber

	return i
}
