package ui

import "github.com/charmbracelet/bubbles/textarea"

func BodyTextarea() textarea.Model {
	i := textarea.New()
	return i
}
