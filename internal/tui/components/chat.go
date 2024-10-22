package components

import (
	"github.com/Broderick-Westrope/teatime/internal/data"
	"github.com/Broderick-Westrope/teatime/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

type ChatModel struct {
	conversation []data.Message
	username     string

	styles    *ChatStyles
	styleFunc func(width, height int) *ChatStyles
}

// NewChatModel creates a new ChatModel.
//   - conversation is list of all messages to display ordered from oldest (first) to newest (last).
//   - username is the username that the active user signed up with. This is used to identify which conversation they have sent.
func NewChatModel(conversation []data.Message, username string) *ChatModel {
	return &ChatModel{
		conversation: conversation,
		username:     username,
		styleFunc:    DefaultChatStyleFunc,
		styles:       DefaultChatStyleFunc(0, 0),
	}
}

func (m *ChatModel) Init() tea.Cmd {
	return nil
}

func (m *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tui.ComponentSizeMsg:
		m.styles = m.styleFunc(msg.Width, msg.Height)
		return m, nil
	}

	return m, nil
}

func (m *ChatModel) View() string {
	var output string

	for _, msg := range m.conversation {
		wasSentByThisUser := msg.Author == m.username
		output += m.viewChatBubble(msg.Content, wasSentByThisUser) + "\n"
	}

	return output
}

func (m *ChatModel) SetConversation(conversation []data.Message) {
	m.conversation = conversation
}

func (m *ChatModel) viewChatBubble(msg string, placeOnRight bool) string {
	var output string

	switch placeOnRight {
	case true:
		output = m.styles.RightBubble.Render(msg)
		output = m.styles.rightAlign.Render(output)

	case false:
		output = m.styles.LeftBubble.Render(msg)
		output = m.styles.leftAlign.Render(output)
	}

	return output
}
