package app

import "github.com/charmbracelet/bubbles/list"

type commandPallete struct {
	commandId, title, desc string
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
	},
	commandPallete{
		commandId: COMMAND_OPEN_BODY,
		title:     "Body",
		desc:      "Edit request body",
	},
	commandPallete{
		commandId: COMMAND_OPEN_HEADER,
		title:     "Header",
		desc:      "Edit request header",
	},
	commandPallete{
		commandId: COMMAND_SELECT_METHOD,
		title:     "Method",
		desc:      "Select request method",
	},
	commandPallete{
		commandId: COMMAND_SAVE_SESSION,
		title:     "Save Session",
		desc:      "Manually save session",
	},
}
