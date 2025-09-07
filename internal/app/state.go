package app

import (
	"github.com/Cybernetics354/mayohttp/internal/app/entity/event"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/intf"
	"github.com/Cybernetics354/mayohttp/internal/app/ui"
	"github.com/Cybernetics354/mayohttp/internal/app/ui/telescope"
	"github.com/Cybernetics354/mayohttp/internal/app/ui/urlcompose"
	"github.com/Cybernetics354/mayohttp/internal/app/val"
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
	resSub       chan event.RequestResult
	pipeResSub   chan event.PipeResult
	commands     list.Model
	sessionList  list.Model
	telescope    telescope.Model
	urlcompose   urlcompose.Model
	url          textinput.Model
	response     textarea.Model
	body         textarea.Model
	header       textarea.Model
	pipe         textinput.Model
	saveInput    textinput.Model
	pipedresp    textarea.Model
	spinner      spinner.Model
	resFilter    intf.ResponseFilter
	showSpinner  bool
	help         help.Model
	keys         val.HomeKeymap
	activity     string
	envList      list.Model
	sw           int
	sh           int
}

func InitialModel() State {
	return State{
		state:        val.STATE_FOCUS_URL,
		stateStack:   []string{val.STATE_FOCUS_URL},
		method:       val.REQUEST_METHOD_GET,
		resSub:       make(chan event.RequestResult),
		pipeResSub:   make(chan event.PipeResult),
		commands:     ui.CommandList(val.CommandPalletes),
		sessionList:  ui.SessionList(),
		telescope:    telescope.New(),
		urlcompose:   urlcompose.New(),
		methodSelect: ui.SelectMethod(val.MethodPalletes),
		body:         ui.BodyTextarea(),
		header:       ui.HeaderTextarea(),
		url:          ui.UrlInput(val.REQUEST_METHOD_GET, ""),
		pipe:         ui.PipeInput(),
		saveInput:    ui.SaveInput(),
		response:     ui.ResponseTextarea(),
		pipedresp:    ui.PipedResponseTextarea(),
		spinner:      ui.Spinner(),
		envList:      ui.EnvList(),
		resFilter:    intf.CreateResponseFilter(),
		showSpinner:  false,
		help:         ui.Help(),
		keys:         val.HomeMapping,
		activity:     "Idle",
		sw:           0,
		sh:           0,
	}
}
