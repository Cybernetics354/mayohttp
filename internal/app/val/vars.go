package val

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	AppStyle = lipgloss.NewStyle().
			Padding(1, 1)
	EnvFilePath        = ".env"
	ConfigFolder       = "./.mayohttp"
	CollectionFolder   = fmt.Sprintf("%s/collections", ConfigFolder)
	TempFilePath       = fmt.Sprintf("%s/temp", ConfigFolder)
	DefaultSessionPath = fmt.Sprintf("%s/session.json", ConfigFolder)
	// debugLogPath       = fmt.Sprintf("%s/debug.log", configFolder)
	ErrorDebugLogPath = fmt.Sprintf("%s/error.log", ConfigFolder)
	HomeLayout        = []string{
		STATE_FOCUS_URL,
		STATE_FOCUS_PIPE,
		STATE_FOCUS_RESPONSE_FILTER,
		STATE_FOCUS_PIPEDRESP,
	}
	Overlays = []string{
		STATE_KEYBINDING_MODAL,
		STATE_SAVE_SESSION_INPUT,
		STATE_SESSION_RENAME_INPUT,
		STATE_TELESCOPE,
		STATE_URL_COMPOSE,
	}
	ResponseSeparator = strings.Repeat("=", 50)
)
