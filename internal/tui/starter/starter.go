package starter

import (
	"github.com/Broderick-Westrope/teatime/internal/tui"
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.child.Update(tui.ComponentSizeMsg{
			Width:  100,
			Height: 20,
		})

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.child, cmd = m.child.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return m.child.View()
}
