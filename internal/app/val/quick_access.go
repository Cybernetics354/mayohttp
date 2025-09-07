package val

import (
	"github.com/Cybernetics354/mayohttp/internal/app/entity/event"
	"github.com/Cybernetics354/mayohttp/internal/app/ui"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var QuickAccess = []list.Item{
	ui.NewListItem(
		"Run Request",
		"Run the whole request from the start",
		[]tea.Msg{
			event.RunRequest{},
		},
	),
	ui.NewListItem(
		"Edit URL",
		"Edit the current session url with default editor",
		[]tea.Msg{
			event.OpenEditor{State: STATE_FOCUS_URL},
		},
	),
	ui.NewListItem("Run Pipe", "Run the pipe only", []tea.Msg{
		event.RunPipe{},
	}),
	ui.NewListItem(
		"Edit Pipe",
		"Edit the current session pipe with default editor",
		[]tea.Msg{
			event.OpenEditor{State: STATE_FOCUS_PIPE},
		},
	),
	ui.NewListItem(
		"Open Response",
		"Open the response in default editor",
		[]tea.Msg{
			event.OpenEditor{State: STATE_FOCUS_PIPEDRESP},
		},
	),
	ui.NewListItem("Edit ENV", "Edit ENV on default editor ($EDITOR)", []tea.Msg{
		event.OpenEnv{},
	}),
	ui.NewListItem("Select ENV", "Select other ENV file", []tea.Msg{
		event.AddStack{State: STATE_SELECT_ENV},
		event.RefreshSelectEnv{},
	}),
}
