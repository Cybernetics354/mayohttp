package app

import "github.com/charmbracelet/bubbles/list"

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
