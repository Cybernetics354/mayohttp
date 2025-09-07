package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Cybernetics354/mayohttp/internal/app"
	comm "github.com/Cybernetics354/mayohttp/internal/app/common"
	"github.com/Cybernetics354/mayohttp/internal/app/val"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	envPtr := flag.String("e", ".env", "Path to the environment file")
	flag.Parse()

	if envPtr != nil && len(*envPtr) >= 0 {
		val.EnvFilePath = *envPtr
	}

	e, err := comm.SetupLogger()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer e.Close()

	p := tea.NewProgram(app.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
