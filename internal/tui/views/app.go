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
	appFocusRegionContacts appFocusRegion = iota
	appFocusRegionChat
)

var _ tea.Model = &AppModel{}

type AppModel struct {
	contacts *components.ConversationsModel
	chat     *components.ChatModel
	focus    appFocusRegion
	styles   *AppStyles
	username string
}

func NewAppModel(conversations []data.Conversation, username string) *AppModel {
	focus := appFocusRegionContacts
	return &AppModel{
		contacts: components.NewConversationsModel(conversations, focus == appFocusRegionContacts),
		chat:     components.NewChatModel(conversations[0], username, focus == appFocusRegionChat),
		focus:    focus,
		styles:   DefaultAppStyles(),
	}
}

func (m *AppModel) Init() tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	cmd = m.contacts.Init()
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
		cmd, err := m.contacts.AddNewMessage(msg.ConversationName, msg.Message)
		if err != nil {
			return m, tui.FatalErrorCmd(err)
		}
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// move from chat to contacts
			if m.focus != appFocusRegionChat {
				break
			}
			err := m.setFocus(appFocusRegionContacts)
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
	case appFocusRegionContacts:
		return tui.UpdateTypedModel(&m.contacts, msg)
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
	case appFocusRegionContacts:
		m.contacts.Enable()
		m.chat.Disable()
	case appFocusRegionChat:
		m.chat.Enable()
		m.contacts.Disable()
	default:
		return fmt.Errorf("unknown appFocusRegion %d", focus)
	}
	m.focus = focus
	return nil
}

// setSize updates calculates and updates the size of the child models taking into account frame sizes.
func (m *AppModel) setSize(windowWidth, windowHeight int) {
	frameWidth, frameHeight := m.styles.TotalFrameSize()
	width, height := windowWidth-frameWidth, windowHeight-frameHeight

	contactsWidth := min(width/3, 35)
	chatWidth := width - contactsWidth

	m.contacts.SetSize(contactsWidth, height)
	m.chat.SetSize(chatWidth, height)
}

func (m *AppModel) View() string {
	output := lipgloss.JoinHorizontal(lipgloss.Center,
		m.styles.Contacts.Render(m.contacts.View()),
		m.styles.Chat.Render(m.chat.View()),
	)
	return m.styles.View.Render(output)
}
