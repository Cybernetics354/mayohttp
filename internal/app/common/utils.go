package common

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Cybernetics354/mayohttp/internal/app/val"
	tea "github.com/charmbracelet/bubbletea"
)

type ParseWithEnvMsg struct {
	Str string
	Err error
}

func SendMsg(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func GetDefaultEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	return editor
}

func ParseWithEnv(str string, c chan ParseWithEnvMsg) {
	command := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(
			"set -a && source %s && set +a && echo '%s' | envsubst",
			val.EnvFilePath,
			strings.ReplaceAll(str, "'", "'\\''"),
		),
	)

	url, err := command.Output()
	if err != nil {
		c <- ParseWithEnvMsg{Str: str, Err: err}
		return
	}

	c <- ParseWithEnvMsg{Str: strings.TrimSpace(string(url)), Err: nil}
}

func Printval(val string, file bool) string {
	// only return as-is for now, i plan to use the bat for pretty printing in near future (or maybe something else if there's better tool)
	if !file {
		return val
	}

	command := exec.Command("cat", val)
	res, err := command.Output()
	if err != nil {
		return ""
	}

	return string(res)
}

func Clamp(val, min, max int) int {
	if val < min {
		return min
	}

	if val > max {
		return max
	}

	return val
}
