package val

import "github.com/charmbracelet/bubbles/list"

type CommandPallete struct {
	CommandId, title, Desc string
}

func (c CommandPallete) Title() string {
	return c.title
}

func (c CommandPallete) Description() string {
	return c.Desc
}

func (c CommandPallete) FilterValue() string {
	return c.title + " " + c.Desc
}

var CommandPalletes = []list.Item{
	CommandPallete{
		CommandId: COMMAND_OPEN_BODY,
		title:     "Body",
		Desc:      "Edit request body",
	},
	CommandPallete{
		CommandId: COMMAND_OPEN_HEADER,
		title:     "Header",
		Desc:      "Edit request header",
	},
	CommandPallete{
		CommandId: COMMAND_SELECT_METHOD,
		title:     "Method",
		Desc:      "Select request method",
	},
	CommandPallete{
		CommandId: COMMAND_OPEN_ENV,
		title:     "ENV",
		Desc:      "Edit ENV file",
	},
	CommandPallete{
		CommandId: COMMAND_CHANGE_ENV,
		title:     "Change ENV",
		Desc:      "Select other ENV file",
	},
	CommandPallete{
		CommandId: COMMAND_OPEN_SESSION_LIST,
		title:     "Open Session",
		Desc:      "Select from collection",
	},
	CommandPallete{
		CommandId: COMMAND_SAVE_SESSION,
		title:     "Save Session",
		Desc:      "Save as collection",
	},
}
