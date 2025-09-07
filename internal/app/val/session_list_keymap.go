package val

import "github.com/charmbracelet/bubbles/key"

type SessionListKeymap struct {
	Up, Down, Filter, Select, Delete, Rename key.Binding
}

func (k SessionListKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Filter, k.Select, k.Rename, k.Delete}
}

func (k SessionListKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var SessionListMapping = SessionListKeymap{
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
		key.WithHelp("enter", "Select"),
	),
	Delete: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "Delete (without confirmation)"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Rename"),
	),
}
