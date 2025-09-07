package ui

import "github.com/charmbracelet/bubbles/textinput"

func UrlInput(method string, initValue string) textinput.Model {
	i := textinput.New()
	i.SetValue(initValue)
	i.Prompt = method + " | "
	i.PromptStyle = UrlPromptStyle

	return i
}
