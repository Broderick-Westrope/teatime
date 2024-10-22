package main

import (
	"log"

	"github.com/Broderick-Westrope/teatime/internal/tui/starter"
	"github.com/Broderick-Westrope/teatime/internal/tui/views"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := starter.NewModel(
		views.NewAppModel(),
	)

	exitModel, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		log.Fatalf("alas, there's been an error: %v\n", err)
	}
	typedExitModel, ok := exitModel.(*starter.Model)
	if !ok {
		log.Fatalln("failed to assert starter model type")
	}
	if typedExitModel.ExitError != nil {
		log.Fatalf("starter model exited with an error: %v\n", typedExitModel.ExitError)
	}
}
