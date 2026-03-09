package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"barbarianprince/ui"
)

func main() {
	forceTutorial := flag.Bool("tutorial", false, "start with in-game tutorial hints enabled")
	flag.Parse()

	model := ui.NewModel(*forceTutorial)
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}