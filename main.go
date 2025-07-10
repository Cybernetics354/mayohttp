package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Cybernetics354/mayohttp/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	envPtr := flag.String("e", ".env", "Path to the environment file")
	flag.Parse()

	if envPtr != nil && len(*envPtr) >= 0 {
		app.EnvFilePath = *envPtr
	}

	p := tea.NewProgram(app.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
