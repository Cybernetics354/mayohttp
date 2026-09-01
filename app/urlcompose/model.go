package urlcompose

import "charm.land/bubbles/v2/textinput"

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
	m.input.SetWidth(m.width)

	return m
}
