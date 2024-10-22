package main

import (
	"log"

	"github.com/Broderick-Westrope/teatime/internal/tui/components"
	"github.com/Broderick-Westrope/teatime/internal/tui/starter"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	msgs := []components.ChatMessage{
		{
			Content: "some",
			Author:  "me",
		},
		{
			Content: "thing",
			Author:  "other",
		},
	}

	m := starter.NewModel(
		components.NewChatModel(msgs, "me", 50),
	)

	_, err := tea.NewProgram(m).Run()
	if err != nil {
		log.Fatal("alas, there's been an error")
	}
}
