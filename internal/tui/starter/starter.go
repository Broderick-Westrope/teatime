package starter

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	child tea.Model
}

func NewModel(child tea.Model) *Model {
	return &Model{
		child: child,
	}
}

func (m *Model) Init() tea.Cmd {
	return m.child.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	m.child, cmd = m.child.Update(msg)

	return m, cmd
}

func (m *Model) View() string {
	return m.child.View()
}
