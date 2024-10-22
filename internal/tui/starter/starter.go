package starter

import (
	"io"

	"github.com/Broderick-Westrope/teatime/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
)

var _ tea.Model = &Model{}

type Model struct {
	child       tea.Model
	ExitError   error
	messagesLog io.Writer
}

func NewModel(child tea.Model, messagesLog io.Writer) *Model {
	return &Model{
		child:       child,
		messagesLog: messagesLog,
	}
}

func (m *Model) Init() tea.Cmd {
	return m.child.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.messagesLog != nil {
		spew.Fdump(m.messagesLog, msg)
	}

	switch msg := msg.(type) {
	case tui.FatalErrorMsg:
		m.ExitError = msg
		return m, tea.Quit

	case tui.SendMessageMsg:
		// TODO: send message to recipient
		var cmd tea.Cmd
		m.child, cmd = m.child.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
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
