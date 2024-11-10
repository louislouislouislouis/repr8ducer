package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/louislouislouislouis/repr8ducer/k8s"
	"github.com/louislouislouislouis/repr8ducer/ui"
)

func main() {
	//// duqmm
	p := tea.NewProgram(
		ui.NewModel(k8s.GetService()),
		tea.WithAltScreen(),
	)

	_, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
