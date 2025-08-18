package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type checkEnvFileMsg struct{}

type openTelescopeMsg struct {
	teleType string
}

type saveSessionMsg struct {
	path string
}

type replaceCurrentSessionMsg struct {
	path string
}

type setupMsg struct{}

type loadSessionListMsg struct{}

type loadSessionMsg struct {
	path string
}

type setFieldValueMsg struct {
	state string
	value string
}

type refreshStateMsg struct{}

type hideSpinnerMsg struct{}

type showSpinnerMsg struct{}

type runRequestMsg struct{}

type nextSectionMsg struct{}

type prevSectionMsg struct{}

type runPipeMsg struct{}

type setActivityMsg string

type openEditorMsg struct {
	state string
}

type openEnvMsg struct{}

type openRequestBodyMsg struct{}

type openRequestHeaderMsg struct{}

type addStackMsg struct {
	state string
}

type popStackMsg struct{}

type popStackRootMsg struct{}

type setStateMsg struct {
	state string
}

type saveInputSubmitMsg struct{}

type selectSessionItemMsg struct{}

type deleteSessionItemMsg struct{}

type selectCommandPalleteMsg struct{}

type selectMethodPalleteMsg struct{}

type runCommandMsg struct {
	commandId string
}

type requestResultMsg struct {
	err error
	res string
}

type pipeResultMsg struct {
	err error
	res string
}

type copyToClipboardMsg struct{}

type refreshSelectEnvMsg struct{}

type selectEnvMsg struct{}

type recalculateComponentSizesMsg struct{}

func listenResponseCmd(sub chan requestResultMsg) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func listenPipeResponseCmd(sub chan pipeResultMsg) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}
