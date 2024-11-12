package starter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"

	"github.com/Broderick-Westrope/teatime/client/internal/db"
	"github.com/Broderick-Westrope/teatime/client/internal/tui"
	"github.com/Broderick-Westrope/teatime/client/internal/tui/views"
	"github.com/Broderick-Westrope/teatime/internal/entity"
	"github.com/Broderick-Westrope/teatime/internal/websocket"
)

var _ tea.Model = &Model{}

type Model struct {
	child       tea.Model
	wsClient    *websocket.Client
	repo        *db.Repository
	messagesLog io.Writer

	creds      *entity.Credentials
	serverAddr string

	width  int
	height int

	msgCh          chan tea.Msg
	cancelWsReader context.CancelFunc
	ExitError      error
}

func NewModel(msgChan chan tea.Msg, serverAddr string, repo *db.Repository, messagesLog io.Writer) (*Model, error) {
	return &Model{
		child:       views.NewLockModel(""),
		repo:        repo,
		messagesLog: messagesLog,
		serverAddr:  serverAddr,
		msgCh:       msgChan,
	}, nil
}

func (m *Model) Init() tea.Cmd {
	return m.child.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.messagesLog != nil {
		_, ok := msg.(tui.AuthenticateMsg)
		if !ok {
			spew.Fdump(m.messagesLog, msg)
		}
	}

	switch msg := msg.(type) {
	case tui.AuthenticateMsg:
		sessionID, err := m.authenticate(context.Background(), msg.IsSignup, msg.Credentials)
		if err != nil {
			if errors.Is(err, errUnauthorised) {
				cmd := m.setChildToLock("Authentication failed, please try again.")
				return m, cmd
			}
			return m, tui.FatalErrorCmd(fmt.Errorf("failed to authenticate: %w", err))
		}
		m.creds = msg.Credentials
		cmd := m.setChildToApp(sessionID)
		return m, cmd

	case tui.QuitMsg:
		switch m.child.(type) {
		case *views.AppModel:
			err := m.appExitCleanup()
			if err != nil {
				return m, tui.FatalErrorCmd(fmt.Errorf("failed to save user data on exit: %w", err))
			}
			cmd := m.setChildToLock("")
			return m, cmd
		default:
			return m, tea.Quit
		}

	case tui.FatalErrorMsg:
		m.ExitError = msg
		_ = m.appExitCleanup()
		return m, tea.Quit

	case tui.SendMessageMsg:
		return m, m.sendMessage(msg.Message, msg.ConversationMD)

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			_ = m.appExitCleanup()
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	m.child, cmd = m.child.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return m.child.View()
}

// sendMessage persists the given message locally and sends it over the relevant WebSocket connections.
// The conversation participants is used to identify which WebSocket clients should receive this message.
func (m *Model) sendMessage(msg entity.Message, conversationMD entity.ConversationMetadata) tea.Cmd {
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
		return tui.FatalErrorCmd(fmt.Errorf("failed to send chat message: %w", err))
	}

	return cmd
}

func (m *Model) authenticate(ctx context.Context, isSignup bool, creds *entity.Credentials) (string, error) {
	route := "/auth/login"
	if isSignup {
		route = "/auth/signup"
	}
	route, err := url.JoinPath(m.serverAddr, route)
	if err != nil {
		return "", fmt.Errorf("failed to join url path: %w", err)
	}

	body, err := json.Marshal(creds)
	if err != nil {
		return "", fmt.Errorf("failed to marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, route, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to build request to %q: %w", route, err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to do request to %q: %w", route, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusCreated:
		break
	case http.StatusUnauthorized:
		return "", errUnauthorised
	default:
		return "", fmt.Errorf("unexpected status code '%d'", resp.StatusCode)
	}

	for _, c := range resp.Cookies() {
		if c.Name == "session_id" {
			return c.Value, nil
		}
	}
	return "", errors.New("failed to find session ID cookie in response")
}

func (m *Model) setChildToLock(errMessage string) tea.Cmd {
	m.child = views.NewLockModel(errMessage)
	cmds := []tea.Cmd{m.child.Init()}

	var cmd tea.Cmd
	m.child, cmd = m.child.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Model) setChildToApp(sessionID string) tea.Cmd {
	wsAddr, err := url.JoinPath(m.serverAddr, "/ws")
	if err != nil {
		return tui.FatalErrorCmd(fmt.Errorf("failed to create WebSocket path: %w", err))
	}
	m.wsClient, err = websocket.NewClient(wsAddr, sessionID)
	if err != nil {
		return tui.FatalErrorCmd(fmt.Errorf("failed to create WebSocket client: %w", err))
	}

	conversations, err := m.repo.GetConversations(m.creds.Username, m.creds.Password)
	if err != nil {
		return tui.FatalErrorCmd(fmt.Errorf("failed to get conversations: %w", err))
	}

	m.child = views.NewAppModel(conversations, m.creds.Username)
	cmds := []tea.Cmd{m.child.Init()}

	var cmd tea.Cmd
	m.child, cmd = m.child.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
	cmds = append(cmds, cmd)

	var ctx context.Context
	ctx, m.cancelWsReader = context.WithCancel(context.Background())
	go m.readFromWebSocket(ctx)

	return tea.Batch(cmds...)
}

func (m *Model) appExitCleanup() error {
	var err error
	// close WS connection
	if m.wsClient != nil {
		err = m.wsClient.Close()
		if err != nil {
			return err
		}
	}
	if m.cancelWsReader != nil {
		m.cancelWsReader()
		m.cancelWsReader = nil
	}

	// save data
	appModel, ok := m.child.(*views.AppModel)
	if !ok {
		return fmt.Errorf("failed to cast starter child to app model: %w", tui.ErrInvalidTypeAssertion)
	}
	conversations, err := appModel.GetConversations()
	if err != nil {
		return err
	}
	err = m.repo.UpdateConversations(m.creds, conversations)
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) readFromWebSocket(ctx context.Context) {
	msgCh := make(chan *websocket.Msg)
	errCh := make(chan error)
	go func() {
		for {
			msgData, err := m.wsClient.ReadMessage()
			if err != nil {
				errCh <- err
				return
			}
			msgCh <- msgData
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return

		case err := <-errCh:
			if websocket.IsNormalCloseError(err) {
				return
			}
			m.msgCh <- tui.FatalErrorMsg(fmt.Errorf("failed to read WebSocket message: %w", err))
			return

		case msg := <-msgCh:
			switch payload := msg.Payload.(type) {
			case websocket.PayloadNotifyConnection:
				m.msgCh <- tui.UpdateConnectionStatusMsg{
					Username:  payload.Username,
					Connected: payload.Connected,
				}
				return

			case websocket.PayloadSendChatMessage:
				if payload.ConversationMD.Name == m.creds.Username {
					payload.ConversationMD.Name = payload.Message.Author
				}
				m.msgCh <- tui.ReceiveMessageMsg{
					ConversationMD: payload.ConversationMD,
					Message:        payload.Message,
				}
				return

			default:
				m.msgCh <- tui.FatalErrorMsg(fmt.Errorf("unknown WebSocket message payload, type=%d", msg.Type))
				return
			}
		}
	}
}
