package val

import (
	"github.com/Cybernetics354/mayohttp/internal/app/ui"
	"github.com/charmbracelet/bubbles/list"
)

var MethodPalletesTelescope = []list.Item{
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

type MethodPallete struct {
	Method string
}

func (c MethodPallete) Title() string {
	return c.Method
}

func (c MethodPallete) Description() string {
	return c.Method
}

func (c MethodPallete) FilterValue() string {
	return c.Method
}

var MethodPalletes = []list.Item{
	MethodPallete{
		Method: REQUEST_METHOD_GET,
	},
	MethodPallete{
		Method: REQUEST_METHOD_POST,
	},
	MethodPallete{
		Method: REQUEST_METHOD_PUT,
	},
	MethodPallete{
		Method: REQUEST_METHOD_DELETE,
	},
	MethodPallete{
		Method: REQUEST_METHOD_OPTIONS,
	},
	MethodPallete{
		Method: REQUEST_METHOD_PATCH,
	},
}
