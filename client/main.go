package main

import (
	"log"
	"os"

	"github.com/Broderick-Westrope/teatime/internal/tui/starter"
	"github.com/Broderick-Westrope/teatime/internal/tui/views"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var messagesLogFile *os.File
	if _, ok := os.LookupEnv("DEBUG"); ok {
		err := os.MkdirAll("logs", 0750)
		if err != nil {
			log.Fatalf("failed to create dir 'logs': %v\n", err)
		}
		messagesLogFile, err = os.OpenFile("logs/messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			log.Fatalf("failed to create file 'logs/messages.log': %v\n", err)
		}
	}

	m := starter.NewModel(
		views.NewAppModel(),
		messagesLogFile,
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
