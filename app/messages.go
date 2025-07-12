package app

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type checkEnvFileMsg struct{}

type saveSessionMsg struct {
	path string
}

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

type setActivityMsg struct {
	activity string
}

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

type setStateMsg struct {
	state string
}

type selectCommandPalleteMsg struct{}

type filterCommandPalleteMsg struct {
	filter list.FilterMatchesMsg
}

type selectMethodPalleteMsg struct{}

type filterMethodPalleteMsg struct {
	filter list.FilterMatchesMsg
}

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

func checkEnvFile() tea.Msg {
	return checkEnvFileMsg{}
}

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

func openEditor(state string) tea.Cmd {
	return func() tea.Msg {
		return openEditorMsg{
			state: state,
		}
	}
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

func filterCommandPallete(msg list.FilterMatchesMsg) tea.Cmd {
	return func() tea.Msg {
		return filterCommandPalleteMsg{filter: msg}
	}
}

func selectMethodPallete() tea.Msg {
	return selectMethodPalleteMsg{}
}

func filterMethodPallete(msg list.FilterMatchesMsg) tea.Cmd {
	return func() tea.Msg {
		return filterMethodPalleteMsg{filter: msg}
	}
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

func saveSession(path string) tea.Cmd {
	return func() tea.Msg {
		return saveSessionMsg{path: path}
	}
}

func loadSession(path string) tea.Cmd {
	return func() tea.Msg {
		return loadSessionMsg{path: path}
	}
}

func openRequestBody() tea.Msg {
	return openRequestBodyMsg{}
}

func openRequestHeader() tea.Msg {
	return openRequestHeaderMsg{}
}

func setActivity(activity string) tea.Cmd {
	return func() tea.Msg {
		return setActivityMsg{activity: activity}
	}
}

func errCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return errMsg(err)
	}
}

func setFieldValue(field string, value string) tea.Msg {
	return setFieldValueMsg{state: field, value: value}
}
