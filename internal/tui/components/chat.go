package components

import (
	"time"

	"github.com/Broderick-Westrope/teatime/internal/data"
	"github.com/Broderick-Westrope/teatime/internal/tui"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = &ChatModel{}

type ChatModel struct {
	conversation []data.Message
	username     string
	contactName  string
	input        textinput.Model
	vp           viewport.Model

	styles    *chatStyles
	styleFunc func(width, height int) *chatStyles
}

// NewChatModel creates a new ChatModel.
//   - conversation is list of all messages to display ordered from oldest (first) to newest (last).
//   - username is the username that the active user signed up with. This is used to identify which conversation they have sent.
func NewChatModel(conversation []data.Message, username, chatName string, enabled bool) *ChatModel {
	input := textinput.New()
	input.Placeholder = "Message"

	styleFunc := disabledChatStyleFunc
	if enabled {
		styleFunc = enabledChatStyleFunc
	}

	return &ChatModel{
		conversation: conversation,
		username:     username,
		contactName:  chatName,
		input:        input,

		styleFunc: styleFunc,
		styles:    styleFunc(0, 0),
	}
}

func (m *ChatModel) Init() tea.Cmd {
	m.switchStyleFunc(m.styleFunc)

	m.vp = viewport.New(0, 0)
	// TODO: derive this position
	m.vp.YPosition = lipgloss.Height(m.viewHeader())
	m.updateViewportContent()

	return tea.Batch(
		m.input.Focus(),
		m.vp.Init(),
	)
}

func (m *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			value := m.input.Value()
			if len(value) == 0 {
				return m, nil
			}
			newMsg := data.Message{
				Content: value,
				Author:  m.username,
				SentAt:  time.Now(),
			}
			m.conversation = append(m.conversation, newMsg)
			m.updateViewportContent()
			m.input.Reset()
			return m, tui.SendMessageCmd(m.contactName, newMsg)
		}
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)
	m.vp, cmd = m.vp.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ChatModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Center,
		m.viewHeader(),
		m.vp.View(),
		m.input.View(),
	)
}

func (m *ChatModel) SetConversation(conversation []data.Message, chatName string) {
	m.conversation = conversation
	m.contactName = chatName
	m.updateViewportContent()
}

func (m *ChatModel) viewHeader() string {
	return m.styles.Header.Render(m.contactName) + "\n"
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
	return m.styles.BubbleStyleFunc(msg, placeOnRight, len(msg)) + "\n"
}

func (m *ChatModel) viewTimestamp(sentAt time.Time) string {
	var output string
	switch {
	case time.Now().Year()-sentAt.Year() > 0: // different year
		output += sentAt.Format("Mon, 02 Jan 2006 at 3:04 AM")

	case time.Now().YearDay()-sentAt.YearDay() == 0: // same day
		output += "Today " + sentAt.Format("3:04 PM")

	case time.Now().YearDay()-sentAt.YearDay() == 1: // previous day
		output += "Yesterday " + sentAt.Format("3:04 PM")

	case time.Now().YearDay()-sentAt.YearDay() < 7: // within the last week
		output += sentAt.Format("Monday 3:04 PM")

	default: // this year but older than a week
		output += sentAt.Format("Mon, 02 Jan at 3:04 PM")
	}
	return m.styles.Timestamp.Render(output) + "\n"
}

func (m *ChatModel) Enable() {
	m.switchStyleFunc(enabledChatStyleFunc)
	m.updateViewportContent()
}

func (m *ChatModel) Disable() {
	m.switchStyleFunc(disabledChatStyleFunc)
	m.updateViewportContent()
}

func (m *ChatModel) switchStyleFunc(styleFunc chatStyleFunc) {
	m.styles = styleFunc(m.styles.Width, m.styles.Height)
	m.styleFunc = styleFunc

	m.input.PromptStyle = m.styles.InputPrompt
	m.input.TextStyle = m.styles.InputText
	m.input.PlaceholderStyle = m.styles.InputPlaceholder
	m.input.CompletionStyle = m.styles.InputCompletion
	m.input.Cursor.Style = m.styles.InputCursor
}

func (m *ChatModel) ResetInput() {
	m.input.Reset()
}

func (m *ChatModel) updateViewportContent() {
	m.vp.SetContent(m.viewConversation())
	m.vp.GotoBottom()
}

func (m *ChatModel) SetSize(width, height int) {
	m.styles = m.styleFunc(width, height)
	m.input.Width = width - (lipgloss.Width(m.input.Prompt) + lipgloss.Width(m.input.Cursor.View()))

	// TODO: derive this margin
	const verticalMargin = 6
	m.vp.Width = width
	m.vp.Height = height - verticalMargin
	m.updateViewportContent()
}
