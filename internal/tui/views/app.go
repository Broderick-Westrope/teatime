package views

import (
	"fmt"

	"github.com/Broderick-Westrope/teatime/internal/entity"
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
	appFocusRegionModal
)

var _ tea.Model = &AppModel{}

type AppModel struct {
	contacts *components.ConversationsModel
	chat     *components.ChatModel
	modal    tea.Model

	focus     appFocusRegion
	prevFocus appFocusRegion
	styles    *AppStyles
	username  string
}

func NewAppModel(conversations []entity.Conversation, username string) *AppModel {
	openConversation := entity.Conversation{
		Name:         "",
		Participants: nil,
		Messages:     nil,
	}
	if len(conversations) > 0 {
		openConversation = conversations[0]
	}

	focus := appFocusRegionContacts
	return &AppModel{
		contacts: components.NewConversationsModel(conversations, focus == appFocusRegionContacts),
		chat:     components.NewChatModel(openConversation, username, focus == appFocusRegionChat),
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

	case tui.OpenModalMsg:
		m.modal = msg.Modal
		m.focus = appFocusRegionModal
		return m, nil

	case tui.SetConversationMsg:
		// setFocus needs to be called before SetConversation since
		// the styling needs to be updated before the viewport is refreshed
		err := m.setFocus(appFocusRegionChat)
		if err != nil {
			return m, tui.FatalErrorCmd(err)
		}
		m.chat.SetConversation(entity.Conversation(msg))
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
			var newFocus appFocusRegion
			switch m.focus {
			case appFocusRegionContacts:
				return m, nil
			case appFocusRegionChat:
				newFocus = appFocusRegionContacts
			case appFocusRegionModal:
				newFocus = m.prevFocus
			}
			err := m.setFocus(newFocus)
			if err != nil {
				return m, tui.FatalErrorCmd(err)
			}
			return m, nil

		case "q":
			if m.focus != appFocusRegionContacts {
				break
			}
			return m, tui.QuitCmd
		}
	}

	cmd, err := m.updateFocussedChild(msg)
	if err != nil {
		return m, tui.FatalErrorCmd(err)
	}
	return m, cmd
}

func (m *AppModel) GetConversations() ([]entity.Conversation, error) {
	return m.contacts.GetConversations()
}

// updateFocussedChild uses the model state to determine which child model to update.
// An error is only returned when an unknown appFocusRegion is provided.
func (m *AppModel) updateFocussedChild(msg tea.Msg) (tea.Cmd, error) {
	switch m.focus {
	case appFocusRegionContacts:
		return tui.UpdateTypedModel(&m.contacts, msg)
	case appFocusRegionChat:
		return tui.UpdateTypedModel(&m.chat, msg)
	case appFocusRegionModal:
		var cmd tea.Cmd
		m.modal, cmd = m.modal.Update(msg)
		return cmd, nil
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
		m.modal = nil
		m.chat.ResetInput()
	case appFocusRegionChat:
		m.chat.Enable()
		m.contacts.Disable()
		m.modal = nil
	case appFocusRegionModal:
		m.chat.Disable()
		m.contacts.Disable()
		m.prevFocus = m.focus
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
	output = m.styles.View.Render(output)

	if m.focus == appFocusRegionModal && m.modal != nil {
		modal := m.styles.Modal.Render(m.modal.View())
		output = tui.OverlayCenter(output, modal)
	}
	return output
}
