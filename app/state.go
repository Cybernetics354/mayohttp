package app

import (
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
	url          textinput.Model
	response     textarea.Model
	body         textarea.Model
	header       textarea.Model
	pipe         textinput.Model
	pipedresp    textarea.Model
	spinner      spinner.Model
	resFilter    ResponseFilter
	showSpinner  bool
	help         help.Model
	keys         keyMap
	activity     string
	err          error
	sw           int
	sh           int
}

func InitialModel() State {
	return State{
		state:        FOCUS_URL,
		stateStack:   []string{FOCUS_URL},
		method:       REQUEST_METHOD_GET,
		resSub:       make(chan requestResultMsg),
		pipeResSub:   make(chan pipeResultMsg),
		commands:     createCommandList(),
		methodSelect: createMethodSelect(),
		body:         createBodyTextarea(),
		header:       createHeaderTextarea(),
		url:          createUrlInput(REQUEST_METHOD_GET),
		pipe:         createPipeInput(),
		response:     createResponseTextarea(),
		pipedresp:    createPipedResponseTextarea(),
		spinner:      createSpinner(),
		resFilter:    CreateResponseFilter(),
		showSpinner:  false,
		help:         createHelp(),
		keys:         keyMaps,
		activity:     "Idle",
		sw:           0,
		sh:           0,
	}
}
