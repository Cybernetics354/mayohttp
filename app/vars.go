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
	configFolder       = "./.mayohttp"
	collectionFolder   = fmt.Sprintf("%s/collections", configFolder)
	tempFilePath       = fmt.Sprintf("%s/temp", configFolder)
	defaultSessionPath = fmt.Sprintf("%s/session.json", configFolder)
	// debugLogPath       = fmt.Sprintf("%s/debug.log", configFolder)
	errorDebugLogPath = fmt.Sprintf("%s/error.log", configFolder)
	homeLayout        = []string{
		STATE_FOCUS_URL,
		STATE_FOCUS_PIPE,
		STATE_FOCUS_RESPONSE_FILTER,
		STATE_FOCUS_PIPEDRESP,
	}
	responseSeparator = strings.Repeat("=", 50)
)
