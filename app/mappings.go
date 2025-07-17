package app

import "github.com/charmbracelet/bubbles/key"

type homeKeymap struct {
	Open, Method, Commands, Quit, Next, Back, Run, Save, SaveAs, OpenSession key.Binding
}

func (k homeKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Open, k.Method, k.Save, k.OpenSession, k.Commands, k.Quit}
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
		key.WithHelp("enter", "Run"),
	),
	Method: key.NewBinding(
		key.WithKeys("ctrl+j"),
		key.WithHelp("ctrl+j", "Select Method"),
	),
	Open: key.NewBinding(
		key.WithKeys("ctrl+o"),
		key.WithHelp("ctrl+o", "Open on editor"),
	),
	Next: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "Next"),
	),
	Back: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "Back"),
	),
	Commands: key.NewBinding(
		key.WithKeys("ctrl+p"),
		key.WithHelp("ctrl+p", "Commands"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "Quit"),
	),
	Save: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "Save"),
	),
	OpenSession: key.NewBinding(
		key.WithKeys("ctrl+l"),
		key.WithHelp("ctrl+l", "Open Session"),
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
		key.WithHelp("enter", "Select"),
	),
}

type saveListKeymap struct {
	Up, Down, Filter, Select, New key.Binding
}

func (k saveListKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Filter, k.Select, k.New}
}

func (k saveListKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var saveListMapping = saveListKeymap{
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

type sessionListKeymap struct {
	Up, Down, Filter, Select, Delete, Rename key.Binding
}

func (k sessionListKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Filter, k.Select, k.Rename, k.Delete}
}

func (k sessionListKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var sessionListMapping = sessionListKeymap{
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
		key.WithHelp("shift+d", "Delete (without confirmation)"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Rename"),
	),
}
