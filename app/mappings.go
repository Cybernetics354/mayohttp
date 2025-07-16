package app

import "github.com/charmbracelet/bubbles/key"

type homeKeymap struct {
	Open, Method, Commands, Quit, Next, Back, Run key.Binding
}

func (k homeKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Open, k.Method, k.Commands, k.Quit}
}

func (k homeKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Commands},
		{k.Quit},
	}
}

var homeMapping = homeKeymap{
	Run: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "Run"),
	),
	Method: key.NewBinding(
		key.WithKeys("ctrl+j"),
		key.WithHelp("Ctrl+j", "Select Method"),
	),
	Open: key.NewBinding(
		key.WithKeys("ctrl+o"),
		key.WithHelp("Ctrl+o", "Open on editor"),
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
		key.WithKeys("ctrl+p"),
		key.WithHelp("Ctrl+p", "Commands"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("Esc", "Quit"),
	),
}

type listKeymap struct {
	Up, Down, Filter, Select key.Binding
}

func (k listKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Filter, k.Select}
}

func (k listKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var listMapping = listKeymap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "Up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "Down"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "Filter"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "Select"),
	),
}
