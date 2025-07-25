package app

import (
	"github.com/Cybernetics354/mayohttp/app/ui"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var quickAccess = []list.Item{
	ui.NewListItem(
		"Run Request",
		"Run the whole request from the start",
		[]tea.Msg{
			runRequestMsg{},
		},
	),
	ui.NewListItem(
		"Edit URL",
		"Edit the current session url with default editor",
		[]tea.Msg{
			openEditorMsg{state: STATE_FOCUS_URL},
		},
	),
	ui.NewListItem("Run Pipe", "Run the pipe only", []tea.Msg{
		runPipeMsg{},
	}),
	ui.NewListItem(
		"Edit Pipe",
		"Edit the current session pipe with default editor",
		[]tea.Msg{
			openEditorMsg{state: STATE_FOCUS_PIPE},
		},
	),
	ui.NewListItem(
		"Open Response",
		"Open the response in default editor",
		[]tea.Msg{
			openEditorMsg{state: STATE_FOCUS_PIPEDRESP},
		},
	),
	ui.NewListItem("Edit ENV", "Edit ENV on default editor ($EDITOR)", []tea.Msg{
		openEnvMsg{},
	}),
	ui.NewListItem("Select ENV", "Select other ENV file", []tea.Msg{
		addStackMsg{state: STATE_SELECT_ENV},
		refreshSelectEnvMsg{},
	}),
}
