package ui

import "charm.land/bubbles/v2/textinput"

func UrlInput(method string, initValue string) textinput.Model {
	i := textinput.New()
	i.SetValue(initValue)
	i.Prompt = method + " | "
	s := i.Styles()
	s.Focused.Prompt = UrlPromptStyle
	s.Blurred.Prompt = UrlPromptStyle
	i.SetStyles(s)

	return i
}
