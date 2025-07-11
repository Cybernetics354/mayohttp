package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type errMsg error

var (
	appStyle = lipgloss.NewStyle().
			Padding(1, 1)
	EnvFilePath        = ".env"
	tempFilePath       = "./.mayohttp/temp"
	debugInitialUrl    = "https://swapi.tech/api/people"
	defaultSessionPath = "./.mayohttp/session.json"
	homeLayout         = []string{FOCUS_URL, FOCUS_PIPE, FOCUS_RESPONSE_FILTER, FOCUS_PIPEDRESP}
	responseSeparator  = fmt.Sprintf("%s", strings.Repeat("=", 50))
)
