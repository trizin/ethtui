package main

import (
	"fmt"
	"os"

	"ethercli/pkg/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	m := ui.GetUI()

	p := tea.NewProgram(m)
	p.EnterAltScreen()

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
