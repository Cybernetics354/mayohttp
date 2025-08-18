package telescope

import (
	"github.com/Cybernetics354/mayohttp/app/ui"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	title          string
	teleType       any
	list           list.Model
	search         textinput.Model
	width          int
	height         int
	searchRegistry map[any]string
}

type SubmitMsg struct {
	Title    string
	TeleType any
	Value    ui.ListItem
}

type ErrorMsg string

func New() Model {
	t := Model{
		title:          "Telescope",
		search:         textinput.New(),
		list:           ui.CreateList(),
		searchRegistry: make(map[any]string),
		width:          60,
		height:         10,
	}

	t.Focus()
	t.RefreshRender()

	return t
}
