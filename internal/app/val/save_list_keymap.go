package val

import "github.com/charmbracelet/bubbles/key"

type SaveListKeymap struct {
	Up, Down, Filter, Select, New key.Binding
}

func (k SaveListKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Filter, k.Select, k.New}
}

func (k SaveListKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var SaveListMapping = SaveListKeymap{
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
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "New"),
	),
}
