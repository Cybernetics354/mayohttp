package app

import "github.com/charmbracelet/lipgloss"

type errMsg error

var (
	appStyle = lipgloss.NewStyle().
			Padding(1, 1)
	tempFilePath    = "./.mayohttp/temp"
	debugInitialUrl = "https://swapi.tech/api/people"
	homeLayout      = []string{FOCUS_URL, FOCUS_PIPE, FOCUS_PIPEDRESP}
)
