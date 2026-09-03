package app

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	"github.com/Cybernetics354/mayohttp/app/component"
	"github.com/Cybernetics354/mayohttp/app/telescope"
	"github.com/Cybernetics354/mayohttp/app/ui"
	"github.com/Cybernetics354/mayohttp/app/urlcompose"
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
	urlcompose   urlcompose.Model
	url          component.UrlInput
	pipe         component.PipeInput
	resFilter    component.ResponseFilter
	pipedResp    viewport.Model
	response     textarea.Model
	body         textarea.Model
	header       textarea.Model
	saveInput    textinput.Model
	spinner      spinner.Model
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
		urlcompose:   urlcompose.New(),
		methodSelect: ui.SelectMethod(methodPalletes),
		body:         ui.BodyTextarea(),
		header:       ui.HeaderTextarea(),
		url:          component.NewUrlInput(REQUEST_METHOD_GET),
		pipe:         component.NewPipeInput(),
		resFilter:    component.NewResponseFilter(),
		saveInput:    ui.SaveInput(),
		response:     ui.ResponseTextarea(),
		pipedResp:    ui.PipedResponseViewport(),
		spinner:      ui.Spinner(),
		envList:      ui.EnvList(),
		showSpinner:  false,
		help:         ui.Help(),
		keys:         homeMapping,
		activity:     "Idle",
		sw:           0,
		sh:           0,
	}
}
