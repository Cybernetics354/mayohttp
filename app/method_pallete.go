package app

import (
	"github.com/Cybernetics354/mayohttp/app/ui"
	"github.com/charmbracelet/bubbles/list"
)

var methodPalletesTelescope = []list.Item{
	ui.NewListItem(
		REQUEST_METHOD_GET,
		"",
		REQUEST_METHOD_GET,
	),
	ui.NewListItem(
		REQUEST_METHOD_POST,
		"",
		REQUEST_METHOD_POST,
	),
	ui.NewListItem(
		REQUEST_METHOD_PUT,
		"",
		REQUEST_METHOD_PUT,
	),
	ui.NewListItem(
		REQUEST_METHOD_DELETE,
		"",
		REQUEST_METHOD_DELETE,
	),
	ui.NewListItem(
		REQUEST_METHOD_OPTIONS,
		"",
		REQUEST_METHOD_OPTIONS,
	),
	ui.NewListItem(
		REQUEST_METHOD_PATCH,
		"",
		REQUEST_METHOD_PATCH,
	),
}

type methodPallete struct {
	method string
}

func (c methodPallete) Title() string {
	return c.method
}

func (c methodPallete) Description() string {
	return c.method
}

func (c methodPallete) FilterValue() string {
	return c.method
}

var methodPalletes = []list.Item{
	methodPallete{
		method: REQUEST_METHOD_GET,
	},
	methodPallete{
		method: REQUEST_METHOD_POST,
	},
	methodPallete{
		method: REQUEST_METHOD_PUT,
	},
	methodPallete{
		method: REQUEST_METHOD_DELETE,
	},
	methodPallete{
		method: REQUEST_METHOD_OPTIONS,
	},
	methodPallete{
		method: REQUEST_METHOD_PATCH,
	},
}
