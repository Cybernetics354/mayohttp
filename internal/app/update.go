package app

import (
	"errors"

	c "github.com/Cybernetics354/mayohttp/internal/app/common"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/event"
	"github.com/Cybernetics354/mayohttp/internal/app/ui/telescope"
	"github.com/Cybernetics354/mayohttp/internal/app/ui/urlcompose"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var command tea.Cmd

	switch msg := msg.(type) {
	case event.Setup:
		return m.Setup()
	case event.CheckEnvFile:
		return m.CheckOrCreateEnvFile()
	case event.SaveSession:
		go m.SaveSession(msg)
	case event.LoadSession:
		return m.LoadSession(msg)
	case event.LoadSessionList:
		return m.LoadSessionList()
	case event.ReplaceCurrentSession:
		return m.ReplaceCurrentSession(msg)
	case list.FilterMatchesMsg:
		return m.HandleListFilter(msg)
	case spinner.TickMsg:
		m.spinner, command = m.spinner.Update(msg)
	case tea.WindowSizeMsg:
		return m.HandleWindowChange(msg)
	case tea.KeyMsg:
		return m.HandleKeyPress(msg)
	case event.RecalculateComponentSizes:
		m.RefreshView()
	case event.SetState:
		return m.SetState(msg.State)
	case event.AddStack:
		return m.AddStack(msg.State)
	case event.PopStack:
		return m.PopStack()
	case event.PopStackRoot:
		return m.PopStackRoot()
	case event.NextSection:
		return m.NextSection()
	case event.PrevSection:
		return m.PrevSection()
	case event.SelectCommandPallete:
		return m.SelectCommandPallete()
	case event.SelectMethodPallete:
		return m.SelectMethodPallete()
	case event.RunCommand:
		return m.RunCommand(msg)
	case event.OpenEnv:
		return m.OpenEnv()
	case event.OpenEditor:
		return m.OpenEditor(msg)
	case event.OpenRequestBody:
		return m.OpenRequestBody()
	case event.OpenRequestHeader:
		return m.OpenRequestHeader()
	case event.HideSpinner:
		return m.HideSpinner()
	case event.ShowSpinner:
		return m.ShowSpinner()
	case event.RefreshState:
		return m.RefreshState()
	case event.RunRequest:
		return m.RunRequest()
	case event.RunPipe:
		return m.RunPipe()
	case event.RequestResult:
		return m.HandleRequestResult(msg)
	case event.PipeResult:
		return m.HandlePipeResult(msg)
	case event.SelectEnv:
		return m.SelectEnv()
	case event.SelectSessionItem:
		return m.SelectSessionItem()
	case event.DeleteSessionItem:
		return m.DeleteSessionItem()
	case event.SaveInputSubmit:
		return m.SaveSessionInputSubmit()
	case event.RefreshSelectEnv:
		return m.RefreshSelectEnv()
	case event.SetActivity:
		m.activity = string(msg)
	case event.SetFieldValue:
		return m.SetFieldValue(msg)
	case event.CopyToClipboard:
		return m.CopyToClipboard()
	case event.OpenTelescope:
		return m.OpenTelescope(msg)
	case telescope.SubmitMsg:
		return m.SelectTelescopeItem(msg)
	case telescope.ErrorMsg:
		return m, c.SendMsg(event.Err(errors.New(string(msg))))
	case urlcompose.Changed:
		m.url.SetValue(msg.Url)
	case urlcompose.Error:
		return m, c.SendMsg(event.Err(errors.New(string(msg))))
	case event.Err:
		if err, ok := msg.(error); ok {
			return m, c.SendMsg(event.Err(err))
		}

		return m.HandleErrorMsg(msg)
	}

	return m, command
}
