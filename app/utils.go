package app

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type parseWithEnvMsg struct {
	str string
	err error
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
	// for now, i plan to use the bat for pretty printing
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
