package views

import (
	"fmt"
	"time"

	"github.com/Broderick-Westrope/charmutils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"

	"github.com/Broderick-Westrope/teatime/client/internal/tui"
	"github.com/Broderick-Westrope/teatime/client/internal/tui/components"
	"github.com/Broderick-Westrope/teatime/internal/entity"
)

// An appFocusRegion is an enum representing a top-level component within the app
// which can become the focus of user navigation. Each value corresponds to a child model.
type appFocusRegion int

const (
	appFocusRegionContacts appFocusRegion = iota
	appFocusRegionChat
	appFocusRegionModal
)

var emptyConversation = entity.Conversation{
	Metadata: entity.ConversationMetadata{
		ID:           uuid.New(),
		Name:         "",
		Participants: nil,
	},
	Messages: nil,
}

var _ tea.Model = &AppModel{}

type AppModel struct {
	conversations *components.ConversationsModel
	chat          *components.ChatModel
	modal         tui.Modal

	focus       appFocusRegion
	prevFocus   appFocusRegion
	styles      *AppStyles
	username    string
	modalWidth  int
	modalHeight int
}

func NewAppModel(conversations []entity.Conversation, username string) *AppModel {
	openConversation := emptyConversation
	if len(conversations) > 0 {
		openConversation = conversations[0]
	}

	focus := appFocusRegionContacts
	return &AppModel{
		conversations: components.NewConversationsModel(conversations, focus == appFocusRegionContacts),
		chat:          components.NewChatModel(openConversation, username, focus == appFocusRegionChat),
		focus:         focus,
		styles:        DefaultAppStyles(),
		username:      username,
	}
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

	case tui.OpenModalMsg:
		cmd := m.setModal(msg.Modal)
		return m, cmd

	case tui.CloseModalMsg:
		cmd := m.setModal(nil)
		return m, cmd

	case tui.CreateConversationMsg:
		cmd := m.createConversation(msg.Name, msg.Participants, msg.NotifyParticipants)
		return m, cmd

	case tui.DeleteConversationMsg:
		cmd := m.deleteConversation(msg.ConversationMD)
		return m, cmd

	case tui.SetConversationMsg:
		cmd := m.setConversation(entity.Conversation(msg))
		return m, cmd

	case tui.ReceiveMessageMsg:
		m.chat.AddNewMessage(msg.Message)
		cmd, err := m.conversations.AddNewMessage(msg.ConversationMD, msg.Message)
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
	return m.conversations.GetConversations()
}

// updateFocussedChild uses the model state to determine which child model to update.
// An error is only returned when an unknown appFocusRegion is provided.
func (m *AppModel) updateFocussedChild(msg tea.Msg) (tea.Cmd, error) {
	switch m.focus {
	case appFocusRegionContacts:
		return charmutils.UpdateTypedModel(&m.conversations, msg)
	case appFocusRegionChat:
		return charmutils.UpdateTypedModel(&m.chat, msg)
	case appFocusRegionModal:
		return charmutils.UpdateTypedModel(&m.modal, msg)
	default:
		return nil, fmt.Errorf("unknown appFocusRegion %d", m.focus)
	}
}

// setFocus enables/disables the child models depending on the provided appFocusRegion and updates the model state.
// An error is only returned when an unknown appFocusRegion is provided.
func (m *AppModel) setFocus(focus appFocusRegion) error {
	switch focus {
	case appFocusRegionContacts:
		m.conversations.Enable()
		m.chat.Disable()
		m.modal = nil
		m.chat.ResetInput()
	case appFocusRegionChat:
		m.chat.Enable()
		m.conversations.Disable()
		m.modal = nil
	case appFocusRegionModal:
		m.modal.Init()
		m.chat.Disable()
		m.conversations.Disable()
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

	m.conversations.SetSize(contactsWidth, height)
	m.chat.SetSize(chatWidth, height)

	m.modalWidth = min(70, width)
	m.modalHeight = height
	if m.modal != nil {
		m.modal.SetSize(m.modalWidth, m.modalHeight)
	}
}

func (m *AppModel) createConversation(name string, participants []string, notifyParticipants bool) tea.Cmd {
	var cmds []tea.Cmd
	conversation := entity.Conversation{
		Metadata: entity.ConversationMetadata{
			ID:           uuid.New(),
			Name:         name,
			Participants: participants,
		},
		Messages: make([]entity.Message, 0),
	}
	cmds = append(cmds, m.conversations.AddNewConversation(conversation))
	m.chat.SetConversation(conversation)

	if notifyParticipants {
		cmd := tui.SendMessageCmd(entity.Message{
			Content: fmt.Sprintf("%q created this conversation 🎉", m.username),
			Author:  m.username,
			SentAt:  time.Now(),
		}, conversation.Metadata)
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}

func (m *AppModel) deleteConversation(conversationMD entity.ConversationMetadata) tea.Cmd {
	if m.chat.GetConversationID() == conversationMD.ID {
		m.chat.SetConversation(entity.Conversation{
			Metadata: entity.ConversationMetadata{
				ID:           uuid.New(),
				Name:         "",
				Participants: nil,
			},
			Messages: nil,
		})
	}
	err := m.conversations.RemoveConversation(conversationMD)
	if err != nil {
		return tui.FatalErrorCmd(err)
	}
	return nil
}

func (m *AppModel) setConversation(conversation entity.Conversation) tea.Cmd {
	// setFocus needs to be called before SetConversation since
	// the styling needs to be updated before the viewport is refreshed
	err := m.setFocus(appFocusRegionChat)
	if err != nil {
		return tui.FatalErrorCmd(err)
	}
	m.chat.SetConversation(conversation)
	return nil
}

func (m *AppModel) setModal(modal tui.Modal) tea.Cmd {
	switch modal == nil {
	case false:
		m.modal = modal
		m.modal.SetSize(m.modalWidth, m.modalHeight)
		err := m.setFocus(appFocusRegionModal)
		if err != nil {
			return tui.FatalErrorCmd(err)
		}
		return nil
	default:
		m.modal = nil
		err := m.setFocus(m.prevFocus)
		if err != nil {
			return tui.FatalErrorCmd(err)
		}
		return nil
	}
}

func (m *AppModel) View() string {
	output := lipgloss.JoinHorizontal(lipgloss.Center,
		m.styles.Contacts.Render(m.conversations.View()),
		m.styles.Chat.Render(m.chat.View()),
	)
	output = m.styles.View.Render(output)

	if m.focus == appFocusRegionModal && m.modal != nil {
		modal := m.styles.Modal.Render(m.modal.View())
		var err error
		output, err = charmutils.OverlayCenter(output, modal, true)
		if err != nil {
			return "** FAILED TO OVERLAY MODAL **" + output
		}
	}
	return output
}
