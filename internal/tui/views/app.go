package views

import tea "github.com/charmbracelet/bubbletea"

type AppModel struct {
}

func NewAppModel() *AppModel {
	return &AppModel{}
}

func (m *AppModel) Init() tea.Cmd {
	return nil
}

func (m *AppModel) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *AppModel) View() string {
	return ""
}
