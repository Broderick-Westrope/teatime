package modals

import (
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = &CreateConversationModel{}

type CreateConversationModel struct{}

func NewCreateConversationModel() *CreateConversationModel {
	return &CreateConversationModel{}
}

func (m *CreateConversationModel) Init() tea.Cmd {
	return nil
}

func (m *CreateConversationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *CreateConversationModel) View() string {
	return "Creating a new conversation...\nPress \"esc\" to go back."
}
