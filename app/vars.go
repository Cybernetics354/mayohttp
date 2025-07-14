package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().
			Padding(1, 1)
	EnvFilePath        = ".env"
	tempFilePath       = "./.mayohttp/temp"
	debugInitialUrl    = "https://swapi.tech/api/people"
	defaultSessionPath = "./.mayohttp/session.json"
	homeLayout         = []string{
		STATE_FOCUS_URL,
		STATE_FOCUS_PIPE,
		STATE_FOCUS_RESPONSE_FILTER,
		STATE_FOCUS_PIPEDRESP,
	}
	responseSeparator = fmt.Sprintf("%s", strings.Repeat("=", 50))
	debugLogPath      = "debug.log"
	errorDebugLogPath = "error.log"
)
