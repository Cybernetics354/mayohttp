package urlcompose

import "github.com/charmbracelet/bubbles/textinput"

type Model struct {
	/// cache the result of the url
	url string

	paths    []string
	queries  map[string]string
	protocol string
	input    textinput.Model

	width int
}

type Error string

type Changed struct {
	Url string
}

func New() Model {
	m := Model{
		url:      "",
		paths:    []string{},
		queries:  make(map[string]string),
		protocol: "",
		input:    textinput.New(),
		width:    60,
	}

	m.input.Focus()
	m.input.ShowSuggestions = true
	m.input.Width = m.width

	m.SetUrl("https://swapi.tech/lorem?walang=gago")

	return m
}
