package components

import (
	"github.com/Broderick-Westrope/teatime/internal/data"
	"github.com/Broderick-Westrope/teatime/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

type ChatModel struct {
	messages []data.Message
	username string
	styles   *ChatStyles
}

// NewChatModel creates a new ChatModel.
//   - messages is list of all messages to display ordered from oldest (first) to newest (last).
//   - username is the username that the active user signed up with. This is used to identify which messages they have sent.
func NewChatModel(messages []data.Message, username string) *ChatModel {
	return &ChatModel{
		messages: messages,
		username: username,
		styles:   DefaultChatStyles(),
	}
}

func (m *ChatModel) Init() tea.Cmd {
	return nil
}

func (m *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tui.ComponentSizeMsg:
		m.styles.leftAlign = m.styles.leftAlign.Width(msg.Width)
		m.styles.rightAlign = m.styles.rightAlign.Width(msg.Width)
		return m, nil
	}

	return m, nil
}

func (m *ChatModel) View() string {
	var output string

	for _, msg := range m.messages {
		wasSentByThisUser := msg.Author == m.username
		output += m.getChatBubble(wasSentByThisUser, msg.Content) + "\n"
	}

	return output
}

func (m *ChatModel) getChatBubble(placeOnRight bool, msg string) string {
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
