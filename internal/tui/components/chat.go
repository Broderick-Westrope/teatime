package components

import (
	"time"

	"github.com/Broderick-Westrope/teatime/internal/data"
	"github.com/Broderick-Westrope/teatime/internal/tui"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ChatModel struct {
	conversation []data.Message
	username     string
	chatName     string
	input        textinput.Model

	styles    *ChatStyles
	styleFunc func(width, height int) *ChatStyles
}

// NewChatModel creates a new ChatModel.
//   - conversation is list of all messages to display ordered from oldest (first) to newest (last).
//   - username is the username that the active user signed up with. This is used to identify which conversation they have sent.
func NewChatModel(conversation []data.Message, username, chatName string) *ChatModel {
	input := textinput.New()
	input.Placeholder = "Message"

	return &ChatModel{
		conversation: conversation,
		username:     username,
		chatName:     chatName,
		input:        input,

		styleFunc: DefaultChatStyleFunc,
		styles:    DefaultChatStyleFunc(0, 0),
	}
}

func (m *ChatModel) Init() tea.Cmd {
	return m.input.Focus()
}

func (m *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tui.ComponentSizeMsg:
		m.styles = m.styleFunc(msg.Width, msg.Height)
		m.input.Width = msg.Width
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			value := m.input.Value()
			m.conversation = append(m.conversation,
				data.Message{
					Content: value,
					Author:  m.username,
					SentAt:  time.Now(),
				},
			)
			m.input.Reset()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *ChatModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Center,
		m.styles.Header.Render(m.chatName)+"\n",
		m.viewConversation(),
		m.input.View(),
	)
}

func (m *ChatModel) SetConversation(conversation []data.Message, chatName string) {
	m.conversation = conversation
	m.chatName = chatName
}

func (m *ChatModel) viewConversation() string {
	var output string
	for i, msg := range m.conversation {
		wasSentByThisUser := msg.Author == m.username
		bubble := m.viewChatBubble(msg.Content, wasSentByThisUser)

		if i == 0 {
			output += m.viewTimestamp(msg.SentAt)
			output += bubble
			continue
		}

		prevMsg := m.conversation[i-1]
		switch {
		case msg.SentAt.Sub(prevMsg.SentAt).Hours() > 12:
			fallthrough
		case msg.SentAt.Sub(prevMsg.SentAt).Hours() > 3 &&
			prevMsg.SentAt.Day() < msg.SentAt.Day():
			output += m.viewTimestamp(msg.SentAt)
		}
		output += bubble
	}
	return m.styles.Conversation.Render(output) + "\n"
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

	return output + "\n"
}

func (m *ChatModel) viewTimestamp(sentAt time.Time) string {
	var output string
	switch {
	case time.Now().Year()-sentAt.Year() > 0: // different year
		output += sentAt.Format("Mon, 02 Jan 2006 at 6:00 AM")

	case time.Now().YearDay()-sentAt.YearDay() == 0: // same day
		output += "Today " + sentAt.Format("6:00 AM")

	case time.Now().YearDay()-sentAt.YearDay() == 1: // previous day
		output += "Yesterday " + sentAt.Format("6:00 AM")

	case time.Now().YearDay()-sentAt.YearDay() < 7: // within the last week
		output += sentAt.Format("Monday 6:00 AM")

	default: // this year but older than a week
		output += sentAt.Format("Mon, 02 Jan at 6:00 AM")
	}
	return m.styles.Timestamp.Render(output) + "\n"
}
