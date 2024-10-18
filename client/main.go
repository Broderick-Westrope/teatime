package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := &Model{}

	_, err := tea.NewProgram(m).Run()
	if err != nil {
		log.Fatal("alas, there's been an error")
	}
}

type Model struct{}

func NewModel() *Model {
	return &Model{}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *Model) View() string {
	return "something"
}
