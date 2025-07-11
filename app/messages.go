package app

import tea "github.com/charmbracelet/bubbletea"

type refreshStateMsg struct{}

type hideSpinnerMsg struct{}

type showSpinnerMsg struct{}

type runRequestMsg struct{}

type nextSectionMsg struct{}

type prevSectionMsg struct{}

type runPipeMsg struct{}

type openEditorMsg struct{}

type openEnvMsg struct{}

type addStackMsg struct {
	state string
}

type popStackMsg struct{}

type setStateMsg struct {
	state string
}

type selectCommandPalleteMsg struct{}

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

type recalculateComponentSizesMsg struct{}

func runRequest() tea.Msg {
	return runRequestMsg{}
}

func runPipe() tea.Msg {
	return runPipeMsg{}
}

func hideSpinner() tea.Msg {
	return hideSpinnerMsg{}
}

func showSpinner() tea.Msg {
	return showSpinnerMsg{}
}

func refreshState() tea.Msg {
	return refreshStateMsg{}
}

func listenResponse(sub chan requestResultMsg) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func listenPipeResponse(sub chan pipeResultMsg) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func openEditor() tea.Msg {
	return openEditorMsg{}
}

func addStack(state string) tea.Cmd {
	return func() tea.Msg {
		return addStackMsg{state}
	}
}

func setState(state string) tea.Cmd {
	return func() tea.Msg {
		return setStateMsg{state}
	}
}

func popStack() tea.Msg {
	return popStackMsg{}
}

func nextSection() tea.Msg {
	return nextSectionMsg{}
}

func prevSection() tea.Msg {
	return prevSectionMsg{}
}

func selectCommandPallete() tea.Msg {
	return selectCommandPalleteMsg{}
}

func openEnv() tea.Msg {
	return openEnvMsg{}
}

func runCommand(commandId string) tea.Cmd {
	return func() tea.Msg {
		return runCommandMsg{commandId: commandId}
	}
}

func recalculateComponentSizes() tea.Msg {
	return recalculateComponentSizesMsg{}
}
