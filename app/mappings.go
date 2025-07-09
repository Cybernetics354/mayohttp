package app

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Commands, Quit, Next, Back, Run key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Commands, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Commands},
		{k.Quit},
	}
}

var mapped = keyMap{
	Run: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "Run"),
	),
	Next: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("Tab", "Next"),
	),
	Back: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("Shift+Tab", "Back"),
	),
	Commands: key.NewBinding(
		key.WithKeys("ctrl+o"),
		key.WithHelp("Ctrl+o", "Commands"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("Esc", "Quit"),
	),
}
