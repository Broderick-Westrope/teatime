package views

import (
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = &LockModel{}

type LockModel struct{}

func NewLockModel() *LockModel {
	return &LockModel{}
}

func (m *LockModel) Init() tea.Cmd {
	//TODO implement me
	panic("implement me")
}

func (m *LockModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//TODO implement me
	panic("implement me")
}

func (m *LockModel) View() string {
	//TODO implement me
	panic("implement me")
}
