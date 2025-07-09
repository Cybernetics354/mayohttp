package app

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func createCommandList() list.Model {
	i := list.New(commandPalletes, list.NewDefaultDelegate(), 0, 0)
	i.Title = "Commands Pallete"
	i.KeyMap.Quit.SetEnabled(false)
	return i
}

func createUrlInput(method string) textinput.Model {
	i := textinput.New()

	i.PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).PaddingRight(2)
	i.Prompt = method
	i.SetValue(debugInitialUrl)
	i.Focus()

	return i
}

func createPipeInput() textinput.Model {
	i := textinput.New()
	i.PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).PaddingRight(2)
	i.Prompt = "PIPE"
	i.Blur()

	return i
}

func createSpinner() spinner.Model {
	i := spinner.New()
	i.Spinner = spinner.Dot
	i.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return i
}

func createBodyTextarea() textarea.Model {
	i := textarea.New()
	return i
}

func createHeaderTextarea() textarea.Model {
	i := textarea.New()
	i.ShowLineNumbers = true
	i.Prompt = ""
	i.FocusedStyle.Base = focusTextarea
	i.BlurredStyle.Base = blurTextarea
	i.Blur()
	return i
}

func createResponseTextarea() textarea.Model {
	i := textarea.New()
	i.ShowLineNumbers = true
	i.Prompt = ""
	i.FocusedStyle.Base = focusTextarea
	i.BlurredStyle.Base = blurTextarea
	i.Blur()
	return i
}

func createPipedResponseTextarea() textarea.Model {
	i := textarea.New()
	i.ShowLineNumbers = true
	i.Prompt = ""
	i.FocusedStyle.Base = focusTextarea
	i.BlurredStyle.Base = blurTextarea
	i.Blur()
	return i
}

func createHelp() help.Model {
	h := help.New()
	return h
}
