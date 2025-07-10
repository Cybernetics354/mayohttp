package app

import "github.com/charmbracelet/bubbles/list"

type commandPallete struct {
	commandId, title, desc, state string
}

func (c commandPallete) Title() string {
	return c.title
}

func (c commandPallete) Description() string {
	return c.desc
}

func (c commandPallete) FilterValue() string {
	return c.title
}

var commandPalletes = []list.Item{
	commandPallete{
		commandId: COMMAND_OPEN_ENV,
		title:     "Open ENV",
		desc:      "Open ENV file with default editor",
		state:     FOCUS_URL,
	},
	commandPallete{
		commandId: COMMAND_OPEN_BODY,
		title:     "Body",
		desc:      "Edit request body",
		state:     FOCUS_URL,
	},
	commandPallete{
		commandId: COMMAND_OPEN_HEADER,
		title:     "Header",
		desc:      "Edit request header",
		state:     FOCUS_URL,
	},
}
