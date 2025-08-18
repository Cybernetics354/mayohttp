package app

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type parseWithEnvMsg struct {
	str string
	err error
}

func sendMsg(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func getDefaultEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	return editor
}

func parseWithEnv(str string, c chan parseWithEnvMsg) {
	command := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(
			"set -a && source %s && set +a && echo '%s' | envsubst",
			EnvFilePath,
			strings.ReplaceAll(str, "'", "'\\''"),
		),
	)

	url, err := command.Output()
	if err != nil {
		c <- parseWithEnvMsg{str: str, err: err}
		return
	}

	c <- parseWithEnvMsg{str: strings.TrimSpace(string(url)), err: nil}
}

func printval(val string, file bool) string {
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

func clamp(val, min, max int) int {
	if val < min {
		return min
	}

	if val > max {
		return max
	}

	return val
}
