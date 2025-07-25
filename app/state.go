package app

import (
	"github.com/Cybernetics354/mayohttp/app/telescope"
	"github.com/Cybernetics354/mayohttp/app/ui"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
)

type State struct {
	state        string
	stateStack   []string
	method       string
	methodSelect list.Model
	resSub       chan requestResultMsg
	pipeResSub   chan pipeResultMsg
	commands     list.Model
	sessionList  list.Model
	telescope    telescope.Model
	url          textinput.Model
	response     textarea.Model
	body         textarea.Model
	header       textarea.Model
	pipe         textinput.Model
	saveInput    textinput.Model
	pipedresp    textarea.Model
	spinner      spinner.Model
	resFilter    ResponseFilter
	showSpinner  bool
	help         help.Model
	keys         homeKeymap
	activity     string
	envList      list.Model
	sw           int
	sh           int
}

func InitialModel() State {
	return State{
		state:        STATE_FOCUS_URL,
		stateStack:   []string{STATE_FOCUS_URL},
		method:       REQUEST_METHOD_GET,
		resSub:       make(chan requestResultMsg),
		pipeResSub:   make(chan pipeResultMsg),
		commands:     ui.CommandList(commandPalletes),
		sessionList:  ui.SessionList(),
		telescope:    telescope.New(),
		methodSelect: ui.SelectMethod(methodPalletes),
		body:         ui.BodyTextarea(),
		header:       ui.HeaderTextarea(),
		url:          ui.UrlInput(REQUEST_METHOD_GET, ""),
		pipe:         ui.PipeInput(),
		saveInput:    ui.SaveInput(),
		response:     ui.ResponseTextarea(),
		pipedresp:    ui.PipedResponseTextarea(),
		spinner:      ui.Spinner(),
		envList:      ui.EnvList(),
		resFilter:    CreateResponseFilter(),
		showSpinner:  false,
		help:         ui.Help(),
		keys:         homeMapping,
		activity:     "Idle",
		sw:           0,
		sh:           0,
	}
}
