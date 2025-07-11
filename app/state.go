package app

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
)

type State struct {
	state       string
	stateStack  []string
	method      string
	resSub      chan requestResponse
	pipeResSub  chan requestPipeResponse
	commands    list.Model
	url         textinput.Model
	response    textarea.Model
	body        textarea.Model
	header      textarea.Model
	pipe        textinput.Model
	pipedresp   textarea.Model
	spinner     spinner.Model
	showSpinner bool
	help        help.Model
	keys        keyMap
	quitting    bool
	err         error
	sw          int
	sh          int
}

func InitialModel() State {
	url := createUrlInput(REQUEST_METHOD_GET)
	url.Focus()

	return State{
		state:       FOCUS_URL,
		stateStack:  []string{FOCUS_URL},
		method:      REQUEST_METHOD_GET,
		resSub:      make(chan requestResponse),
		pipeResSub:  make(chan requestPipeResponse),
		commands:    createCommandList(),
		body:        createBodyTextarea(),
		header:      createHeaderTextarea(),
		url:         url,
		pipe:        createPipeInput(),
		response:    createResponseTextarea(),
		pipedresp:   createPipedResponseTextarea(),
		spinner:     createSpinner(),
		showSpinner: false,
		help:        createHelp(),
		keys:        keyMaps,
		sw:          0,
		sh:          0,
	}
}
