package components

import tea "github.com/charmbracelet/bubbletea"

type ContactsModel struct {
}

func NewContactsModel() *ContactsModel {
	return &ContactsModel{}
}

func (m *ContactsModel) Init() tea.Cmd {
	return nil
}

func (m *ContactsModel) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *ContactsModel) View() string {
	return ""
}
