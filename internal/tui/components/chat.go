package components

import (
	"github.com/Broderick-Westrope/teatime/internal/data"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ChatModel struct {
	messages []data.Message
	username string
	styles   ChatStyles
}

type ChatStyles struct {
	LeftBubble  lipgloss.Style
	RightBubble lipgloss.Style
	leftAlign   lipgloss.Style
	rightAlign  lipgloss.Style
}

// NewChatModel creates a new ChatModel.
//   - messages is list of all messages to display ordered from oldest (first) to newest (last).
//   - username is the username that the active user signed up with. This is used to identify which messages they have sent.
//   - width is the width to use for each chat bubble. It should be equal to the total width of the model. This allows aligning messages on the right.
func NewChatModel(messages []data.Message, username string, width int) *ChatModel {
	leftBubbleBorder := lipgloss.RoundedBorder()
	leftBubbleBorder.BottomLeft = "└"

	rightBubbleBorder := lipgloss.RoundedBorder()
	rightBubbleBorder.BottomRight = "┘"

	return &ChatModel{
		messages: messages,
		username: username,
		styles: ChatStyles{
			LeftBubble:  lipgloss.NewStyle().Border(leftBubbleBorder, true),
			RightBubble: lipgloss.NewStyle().Border(rightBubbleBorder, true),
			leftAlign:   lipgloss.NewStyle().Width(width).AlignHorizontal(lipgloss.Left),
			rightAlign:  lipgloss.NewStyle().Width(width).AlignHorizontal(lipgloss.Right),
		},
	}
}

func (m *ChatModel) Init() tea.Cmd {
	return nil
}

func (m *ChatModel) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
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
