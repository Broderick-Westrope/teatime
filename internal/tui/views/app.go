package views

import (
	"fmt"

	"github.com/Broderick-Westrope/teatime/internal/data"
	"github.com/Broderick-Westrope/teatime/internal/tui"
	"github.com/Broderick-Westrope/teatime/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// An appFocusRegion is an enum representing a top-level component within the app
// which can become the focus of user navigation. Each value corresponds to a child model.
type appFocusRegion int

const (
	appFocusRegionConversations appFocusRegion = iota
	appFocusRegionChat
)

var _ tea.Model = &AppModel{}

type AppModel struct {
	conversations *components.ConversationsModel
	chat          *components.ChatModel
	help          *components.ContextualHelp[appFocusRegion]

	focus   appFocusRegion
	keymaps map[appFocusRegion]components.KeyMap
	styles  *AppStyles
}

func NewAppModel(conversations []data.Conversation, username string) (*AppModel, error) {
	keymaps := defaultAppKeyMaps()
	conversationsKeymap, ok := keymaps[appFocusRegionConversations].(*conversationsKeyMap)
	if !ok {
		return nil, fmt.Errorf("failed to get conversations keymap: %w", tui.ErrInvalidTypeAssertion)
	}

	focus := appFocusRegionConversations

	return &AppModel{
		conversations: components.NewConversationsModel(conversations, conversationsKeymap.ListDelegate, focus == appFocusRegionConversations),
		chat:          components.NewChatModel(conversations[0], username, focus == appFocusRegionChat),
		help:          components.NewContextualHelp(keymaps, focus),

		focus:  focus,
		styles: DefaultAppStyles(),
	}, nil
}

func (m *AppModel) Init() tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	cmd = m.conversations.Init()
	cmds = append(cmds, cmd)

	cmd = m.chat.Init()
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
		return m, nil

	case tui.SetConversationMsg:
		// setFocus needs to be called before SetConversation since
		// the styling needs to be updated before the viewport is refreshed
		err := m.setFocus(appFocusRegionChat)
		if err != nil {
			return m, tui.FatalErrorCmd(err)
		}
		m.chat.SetConversation(data.Conversation(msg))
		return m, nil

	case tui.ReceiveMessageMsg:
		m.chat.AddNewMessage(msg.Message)
		cmd, err := m.conversations.AddNewMessage(msg.ConversationName, msg.Message)
		if err != nil {
			return m, tui.FatalErrorCmd(err)
		}
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "?":
			m.help.ShowAll = !m.help.ShowAll
			return m, nil

		case "esc":
			// move from chat to conversations
			if m.focus != appFocusRegionChat {
				break
			}
			err := m.setFocus(appFocusRegionConversations)
			if err != nil {
				return m, tui.FatalErrorCmd(err)
			}
			m.chat.ResetInput()
			return m, nil

		case "q":
			// don't quit if the user types "q" in the chat
			if m.focus == appFocusRegionChat {
				break
			}
			return m, tea.Quit
		}
	}

	cmd, err := m.updateFocussedChild(msg)
	if err != nil {
		return m, tui.FatalErrorCmd(err)
	}
	return m, cmd
}

// updateFocussedChild uses the model state to determine which child model to update.
// An error is only returned when an unknown appFocusRegion is provided.
func (m *AppModel) updateFocussedChild(msg tea.Msg) (tea.Cmd, error) {
	switch m.focus {
	case appFocusRegionConversations:
		return tui.UpdateTypedModel(&m.conversations, msg)
	case appFocusRegionChat:
		return tui.UpdateTypedModel(&m.chat, msg)
	default:
		return nil, fmt.Errorf("unknown appFocusRegion %d", m.focus)
	}
}

// setFocus enables/disables the child models depending on the provided appFocusRegion and updates the model state.
// An error is only returned when an unknown appFocusRegion is provided.
func (m *AppModel) setFocus(focus appFocusRegion) error {
	switch focus {
	case appFocusRegionConversations:
		m.conversations.Enable()
		m.chat.Disable()
	case appFocusRegionChat:
		m.chat.Enable()
		m.conversations.Disable()
	default:
		return fmt.Errorf("unknown appFocusRegion %d", focus)
	}
	m.focus = focus
	m.help.Current = focus
	return nil
}

// setSize updates calculates and updates the size of the child models taking into account frame sizes.
func (m *AppModel) setSize(windowWidth, windowHeight int) {
	frameWidth, frameHeight := m.styles.TotalFrameSize()
	width, height := windowWidth-frameWidth, windowHeight-frameHeight

	conversationsWidth := min(width/3, 35)
	chatWidth := width - conversationsWidth

	conversationsHeight := height - lipgloss.Height(m.help.View())

	m.conversations.SetSize(conversationsWidth, conversationsHeight)
	m.chat.SetSize(chatWidth, height)
}

func (m *AppModel) View() string {
	output := lipgloss.JoinHorizontal(lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Left,
			m.styles.Conversations.Render(m.conversations.View()),
			m.help.View(),
		),
		m.styles.Chat.Render(m.chat.View()),
	)
	return m.styles.View.Render(output)
}
