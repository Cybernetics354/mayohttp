package app

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var command tea.Cmd

	switch msg := msg.(type) {
	case setupMsg:
		return m.Setup()
	case checkEnvFileMsg:
		return m.CheckOrCreateEnvFile()
	case saveSessionMsg:
		go m.SaveSession(msg)
	case loadSessionMsg:
		return m.LoadSession(msg)
	case loadSessionListMsg:
		return m.LoadSessionList()
	case replaceCurrentSessionMsg:
		return m.ReplaceCurrentSession(msg)
	case list.FilterMatchesMsg:
		return m.HandleListFilter(msg)
	case spinner.TickMsg:
		m.spinner, command = m.spinner.Update(msg)
	case tea.WindowSizeMsg:
		return m.HandleWindowChange(msg)
	case tea.KeyMsg:
		return m.HandleKeyPress(msg)
	case recalculateComponentSizesMsg:
		m.RefreshView()
	case setStateMsg:
		return m.SetState(msg.state)
	case addStackMsg:
		return m.AddStack(msg.state)
	case popStackMsg:
		return m.PopStack()
	case popStackRootMsg:
		return m.PopStackRoot()
	case nextSectionMsg:
		return m.NextSection()
	case prevSectionMsg:
		return m.PrevSection()
	case selectCommandPalleteMsg:
		return m.SelectCommandPallete()
	case selectMethodPalleteMsg:
		return m.SelectMethodPallete()
	case runCommandMsg:
		return m.RunCommand(msg)
	case openEnvMsg:
		return m.OpenEnv()
	case openEditorMsg:
		return m.OpenEditor(msg)
	case openRequestBodyMsg:
		return m.OpenRequestBody()
	case openRequestHeaderMsg:
		return m.OpenRequestHeader()
	case hideSpinnerMsg:
		return m.HideSpinner()
	case showSpinnerMsg:
		return m.ShowSpinner()
	case refreshStateMsg:
		return m.RefreshState()
	case runRequestMsg:
		return m.RunRequest()
	case runPipeMsg:
		return m.RunPipe()
	case requestResultMsg:
		return m.HandleRequestResult(msg)
	case pipeResultMsg:
		return m.HandlePipeResult(msg)
	case selectEnvMsg:
		return m.SelectEnv()
	case selectSessionItemMsg:
		return m.SelectSessionItem()
	case refreshSelectEnvMsg:
		return m.RefreshSelectEnv()
	case setActivityMsg:
		m.activity = string(msg)
	case setFieldValueMsg:
		return m.SetFieldValue(msg)
	case errMsg:
		return m.HandleErrorMsg(msg)
	}

	return m, command
}
