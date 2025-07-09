package app

import "github.com/charmbracelet/bubbles/list"

type commandPallete struct {
	title, desc, state string
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
		title: "URL",
		desc:  "Edit URL",
		state: FOCUS_URL,
	},
	commandPallete{
		title: "Response",
		desc:  "Open response on bigger view",
		state: FOCUS_RESPONSE,
	},
	commandPallete{
		title: "Pipe",
		desc:  "Pipe the response",
		state: FOCUS_PIPE,
	},
	commandPallete{
		title: "Piped Response",
		desc:  "Open piped response on bigger view",
		state: FOCUS_PIPEDRESP,
	},
	commandPallete{
		title: "Body",
		desc:  "Edit request body",
		state: EDIT_BODY,
	},
	commandPallete{
		title: "Header",
		desc:  "Edit request header",
		state: EDIT_HEADER,
	},
}
