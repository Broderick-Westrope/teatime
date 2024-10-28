package components

import (
	"strings"
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
	conversation data.Conversation
	username     string
	input        textinput.Model
	vp           viewport.Model

	styles    *chatStyles
	styleFunc func(width, height int) *chatStyles
}

// NewChatModel creates a new ChatModel.
//   - messages: list of all messages to display ordered from oldest (first) to newest (last).
//   - username: the username that the active user signed up with. This is used to identify which messages they have sent.
//   - chatName: the title to display in the chat header.
//   - enabled: whether this component is enabled to begin with.
func NewChatModel(conversation data.Conversation, username string, enabled bool) *ChatModel {
	input := textinput.New()
	input.Placeholder = "Message"

	styleFunc := disabledChatStyleFunc
	if enabled {
		styleFunc = enabledChatStyleFunc
	}

	return &ChatModel{
		conversation: conversation,
		username:     username,
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
	m.refreshViewportContent()

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
			if len(strings.TrimSpace(value)) == 0 {
				return m, nil
			}
			newMsg := data.Message{
				Content: value,
				Author:  m.username,
				SentAt:  time.Now(),
			}
			m.conversation.Messages = append(m.conversation.Messages, newMsg)
			m.refreshViewportContent()
			m.input.Reset()
			return m, tui.SendMessageCmd(newMsg, m.conversation)
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

// SetConversation updates the model state to have the given messages and chatName.
// It also refreshes the viewport content.
func (m *ChatModel) SetConversation(conversation data.Conversation) {
	m.conversation = conversation
	m.refreshViewportContent()
}

// SetSize calculates and applies the correct size to its nested components.
// The given width and height should be the dimensions for this component not the window.
func (m *ChatModel) SetSize(width, height int) {
	m.styles = m.styleFunc(width, height)
	m.input.Width = width - (lipgloss.Width(m.input.Prompt) + lipgloss.Width(m.input.Cursor.View()))

	// TODO: derive this margin
	const verticalMargin = 6
	m.vp.Width = width
	m.vp.Height = height - verticalMargin
	m.refreshViewportContent()
}

// Enable makes the model appear as though it is active/focussed.
func (m *ChatModel) Enable() {
	m.switchStyleFunc(enabledChatStyleFunc)
	m.refreshViewportContent()
}

// Disable makes the model appear as though it is not active/focussed.
func (m *ChatModel) Disable() {
	m.switchStyleFunc(disabledChatStyleFunc)
	m.refreshViewportContent()
}

// ResetInput will reset the content of the nested input component.
func (m *ChatModel) ResetInput() {
	m.input.Reset()
}

// switchStyleFunc updates the model state to have the given styleFunc.
// It also updates the styles of this component and nested components using the provided styleFunc.
func (m *ChatModel) switchStyleFunc(styleFunc chatStyleFunc) {
	m.styles = styleFunc(m.styles.Width, m.styles.Height)
	m.styleFunc = styleFunc

	m.input.PromptStyle = m.styles.InputPrompt
	m.input.TextStyle = m.styles.InputText
	m.input.PlaceholderStyle = m.styles.InputPlaceholder
	m.input.CompletionStyle = m.styles.InputCompletion
	m.input.Cursor.Style = m.styles.InputCursor
}

// refreshViewportContent recalculates the messages and sets the viewport content to the result.
// It also scrolls the viewport to the bottom to mimick the behaviour of common messaging apps. This should
// be done after updating any styles since it will update the viewport content with the styled messages.
func (m *ChatModel) refreshViewportContent() {
	m.vp.SetContent(m.viewConversation())
	m.vp.GotoBottom()
}

func (m *ChatModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Center,
		m.viewHeader(),
		m.vp.View(),
		m.input.View(),
	)
}

// viewHeader returns the styled output for the header.
func (m *ChatModel) viewHeader() string {
	return m.styles.Header.Render(m.conversation.Name) + "\n"
}

// viewConversation returns the styled output for the messages.
func (m *ChatModel) viewConversation() string {
	var output string
	for i, msg := range m.conversation.Messages {
		wasSentByThisUser := msg.Author == m.username
		bubble := m.viewChatBubble(msg.Content, wasSentByThisUser)

		if i == 0 {
			output += m.viewTimestamp(msg.SentAt)
			output += bubble
			continue
		}

		prevMsg := m.conversation.Messages[i-1]
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

// viewChatBubble returns the styled output for a single chat bubble.
func (m *ChatModel) viewChatBubble(msg string, placeOnRight bool) string {
	return m.styles.BubbleStyleFunc(msg, placeOnRight, len(msg)) + "\n"
}

// viewTimestamp returns the styled output for a single timestamp value.
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
