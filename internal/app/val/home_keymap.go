package val

import "github.com/charmbracelet/bubbles/key"

type HomeKeymap struct {
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
	QuickAccess,
	ComposeUrl,
	ClearInput,
	Keybinding key.Binding
}

func (k HomeKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Run, k.Keybinding}
}

func (k HomeKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Commands},
		{k.Quit},
	}
}

func (k *HomeKeymap) KeybindingHelp() []key.Binding {
	return []key.Binding{
		k.Run,
		k.ComposeUrl,
		k.QuickAccess,
		k.CopyToClipboard,
		k.ClearInput,
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

var HomeMapping = HomeKeymap{
	ClearInput: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("<c-d>", "Clear focused input"),
	),
	ComposeUrl: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("<c-u>", "Open URL composer"),
	),
	QuickAccess: key.NewBinding(
		key.WithKeys("ctrl+j"),
		key.WithHelp("<c-j>", "Open quick access menu"),
	),
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
		key.WithKeys("ctrl+k"),
		key.WithHelp("<c-k>", "Select Request Method (GET, POST, PUT, DELETE, OPTIONS, PATCH)"),
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
