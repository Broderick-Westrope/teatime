package starter

import (
	"fmt"
	"io"

	"github.com/Broderick-Westrope/teatime/internal/db"
	"github.com/Broderick-Westrope/teatime/internal/entity"
	"github.com/Broderick-Westrope/teatime/internal/tui"
	"github.com/Broderick-Westrope/teatime/internal/tui/views"
	"github.com/Broderick-Westrope/teatime/internal/websocket"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
)

var _ tea.Model = &Model{}

type Model struct {
	child       tea.Model
	wsClient    *websocket.Client
	repo        *db.Repository
	messagesLog io.Writer

	username string
	password string

	ExitError error
}

func NewModel(username, password string, wsClient *websocket.Client, repo *db.Repository, messagesLog io.Writer) (*Model, error) {
	conversations, err := repo.GetConversations(username, password)
	if err != nil {
		return nil, err
	}

	return &Model{
		child:       views.NewAppModel(conversations, username),
		wsClient:    wsClient,
		repo:        repo,
		messagesLog: messagesLog,

		username: username,
		password: password,
	}, nil
}

func (m *Model) Init() tea.Cmd {
	return m.child.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.messagesLog != nil {
		spew.Fdump(m.messagesLog, msg)
	}

	switch msg := msg.(type) {
	case tui.QuitMsg:
		err := m.saveUserData()
		if err != nil {
			return m, tui.FatalErrorCmd(fmt.Errorf("failed to save user data on exit: %w", err))
		}
		return m, tea.Quit

	case tui.FatalErrorMsg:
		m.ExitError = msg
		return m, tui.QuitCmd

	case tui.SendMessageMsg:
		return m, m.SendMessage(msg.Message, msg.ConversationMD)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tui.QuitCmd
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
func (m *Model) SendMessage(msg entity.Message, conversationMD entity.ConversationMetadata) tea.Cmd {
	// Add message locally
	var cmd tea.Cmd
	m.child, cmd = m.child.Update(tui.ReceiveMessageMsg{
		ConversationMD: conversationMD,
		Message:        msg,
	})

	// Send message to recipients via WebSockets
	recipients := conversationMD.Participants
	for i, v := range recipients {
		if v == msg.Author {
			recipients = append(recipients[:i], recipients[i+1:]...)
			break
		}
	}

	err := m.wsClient.SendChatMessage(msg, conversationMD, recipients)
	if err != nil {
		return tui.FatalErrorCmd(fmt.Errorf("failed to send chat message: %v\n", err))
	}

	return cmd
}

func (m *Model) saveUserData() error {
	appModel, ok := m.child.(*views.AppModel)
	if !ok {
		return fmt.Errorf("failed to cast starter child to app model: %w", tui.ErrInvalidTypeAssertion)
	}
	conversations, err := appModel.GetConversations()
	if err != nil {
		return err
	}
	err = m.repo.SetConversations(m.username, m.password, conversations)
	return err
}
