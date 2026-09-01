package ui

import "charm.land/bubbles/v2/textarea"

func ResponseTextarea() textarea.Model {
	i := textarea.New()
	i.ShowLineNumbers = true
	i.Prompt = ""
	styles := i.Styles()
	styles.Focused.Base = FocusTextarea
	styles.Blurred.Base = BlurTextarea
	styles.Focused.LineNumber = FocusTextareaLineNumber
	styles.Blurred.LineNumber = BlurTextareaLineNumber
	i.SetStyles(styles)

	return i
}
