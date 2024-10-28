package starter

import (
	"fmt"
	"io"

	"github.com/Broderick-Westrope/teatime/internal/data"
	"github.com/Broderick-Westrope/teatime/internal/tui"
	"github.com/Broderick-Westrope/teatime/internal/websocket"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
)

var _ tea.Model = &Model{}

type Model struct {
	child       tea.Model
	wsClient    *websocket.Client
	messagesLog io.Writer

	ExitError error
}

func NewModel(child tea.Model, wsClient *websocket.Client, messagesLog io.Writer) *Model {
	return &Model{
		child:       child,
		wsClient:    wsClient,
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
		return m, m.SendMessage(msg.Message, msg.Conversation)

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

// SendMessage persists the given message locally and sends it over the relevant WebSocket connections.
// The conversation participants is used to identify which WebSocket clients should receive this message.
func (m *Model) SendMessage(msg data.Message, conversation data.Conversation) tea.Cmd {
	// Add message locally
	var cmd tea.Cmd
	m.child, cmd = m.child.Update(tui.ReceiveMessageMsg{
		ConversationName: conversation.Name,
		Message:          msg,
	})

	// Send message to recipients via WebSockets
	recipients := conversation.Participants
	for i, v := range recipients {
		if v == msg.Author {
			recipients = append(recipients[:i], recipients[i+1:]...)
			break
		}
	}

	conversationName := conversation.Name
	if len(recipients) == 1 {
		conversationName = msg.Author
	}

	err := m.wsClient.SendChatMessage(msg, conversationName, recipients)
	if err != nil {
		return tui.FatalErrorCmd(fmt.Errorf("failed to send chat message: %v\n", err))
	}

	return cmd
}
