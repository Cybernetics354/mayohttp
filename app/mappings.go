package app

import "github.com/charmbracelet/bubbles/key"

type homeKeymap struct {
	Open,
	Method,
	Commands,
	Quit,
	Next,
	Back,
	Run,
	Save,
	SaveAs,
	OpenSession,
	CopyToClipboard,
	OpenEnv,
	Keybinding key.Binding
}

func (k homeKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Run, k.Keybinding}
}

func (k homeKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Commands},
		{k.Quit},
	}
}

func (k *homeKeymap) KeybindingHelp() []key.Binding {
	return []key.Binding{
		k.Run,
		k.CopyToClipboard,
		k.Method,
		k.Open,
		k.Next,
		k.Back,
		k.Commands,
		k.OpenEnv,
		k.Save,
		k.OpenSession,
		k.Quit,
		k.Keybinding,
	}
}

var homeMapping = homeKeymap{
	OpenEnv: key.NewBinding(
		key.WithKeys("ctrl+e"),
		key.WithHelp("<c-e>", "Open ENV on default editor ($EDITOR)"),
	),
	CopyToClipboard: key.NewBinding(
		key.WithKeys("ctrl+y"),
		key.WithHelp("<c-y>", "Copy current focused text field to clipboard"),
	),
	Run: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Run"),
	),
	Method: key.NewBinding(
		key.WithKeys("ctrl+j"),
		key.WithHelp("<c-j>", "Select Request Method (GET, POST, PUT, DELETE, OPTIONS, PATCH)"),
	),
	Open: key.NewBinding(
		key.WithKeys("ctrl+o"),
		key.WithHelp("<c-o>", "Open the current active field on default editor ($EDITOR)"),
	),
	Next: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "Move to next section"),
	),
	Back: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "Move to previous section"),
	),
	Commands: key.NewBinding(
		key.WithKeys("ctrl+p"),
		key.WithHelp("<c-p>", "Open Command List"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "Quit the MayoHTTP"),
	),
	Save: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("<c-s>", "Save current Session"),
	),
	OpenSession: key.NewBinding(
		key.WithKeys("ctrl+l"),
		key.WithHelp("<c-l>", "Open/Load another Session"),
	),
	Keybinding: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("<c-h>", "Keybindings"),
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
		key.WithHelp("D", "Delete (without confirmation)"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Rename"),
	),
}
