package components

import (
	"fmt"

	"github.com/Broderick-Westrope/teatime/client/internal/tui"
	"github.com/Broderick-Westrope/teatime/internal/entity"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = &ConversationsModel{}

type ConversationsModel struct {
	list   list.Model
	styles *conversationsStyles
}

func NewConversationsModel(conversations []entity.Conversation, enabled bool) *ConversationsModel {
	var items = make([]list.Item, len(conversations))
	for i, d := range conversations {
		items[i] = Conversation(d)
	}

	styles := disabledConversationsStyles()
	if enabled {
		styles = enabledConversationsStyles()
	}

	delegate := NewListDelegate(DefaultListDelegateKeyMap(), styles.ListItem)
	conversationsList := list.New(items, delegate, 0, 0)
	conversationsList.Title = "Conversations"
	conversationsList.DisableQuitKeybindings()

	return &ConversationsModel{
		list:   conversationsList,
		styles: styles,
	}
}

func (m *ConversationsModel) Init() tea.Cmd {
	m.switchStyles(m.styles)
	return nil
}

func (m *ConversationsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		// all key presses when filtering should go to the nested list component
		if m.list.FilterState() == list.Filtering {
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *ConversationsModel) AddNewConversation(conversation entity.Conversation) tea.Cmd {
	return m.list.InsertItem(0, Conversation(conversation))
}

// AddNewMessage will add the given message to the chat with the given chatName.
// It will also move this messages to the top of the contacts list and update the list selection.
func (m *ConversationsModel) AddNewMessage(conversationMD entity.ConversationMetadata, message entity.Message) (tea.Cmd, error) {
	foundIdx := -1
	items := m.list.Items()

	for i, item := range items {
		conversation, ok := item.(Conversation)
		if !ok {
			return nil, fmt.Errorf("failed to add new message to conversation: (list item) %v", tui.ErrInvalidTypeAssertion)
		}

		if conversation.Metadata.ID != conversationMD.ID {
			continue
		}

		foundIdx = i
		conversation.Messages = append(conversation.Messages, message)
		items[i] = conversation
		break
	}

	switch {
	case foundIdx < 0:
		item := Conversation{
			Metadata: conversationMD,
			Messages: []entity.Message{message},
		}
		items = append([]list.Item{item}, items...)
		m.list.Select(0)
		return m.list.SetItems(items), nil

	case foundIdx == 0: // already first
		return m.list.SetItems(items), nil

	default: // if the conversation is not already first in the list, move it to first
		item := items[foundIdx]
		items = append(items[:foundIdx], items[foundIdx+1:]...)
		items = append([]list.Item{item}, items...)
		m.list.Select(0)
		return m.list.SetItems(items), nil
	}
}

func (m *ConversationsModel) RemoveConversation(conversationMD entity.ConversationMetadata) error {
	items := m.list.Items()

	for i, item := range items {
		conversation, ok := item.(Conversation)
		if !ok {
			return fmt.Errorf("failed to remove conversation: (list item) %v", tui.ErrInvalidTypeAssertion)
		}

		if conversation.Metadata.ID == conversationMD.ID {
			m.list.RemoveItem(i)
			break
		}
	}
	return nil
}

func (m *ConversationsModel) GetConversations() ([]entity.Conversation, error) {
	items := m.list.Items()
	conversations := make([]entity.Conversation, len(items))
	for i := range items {
		conversation, ok := items[i].(Conversation)
		if !ok {
			return nil, fmt.Errorf("failed to get conversations: %w", tui.ErrInvalidTypeAssertion)
		}
		conversations[i] = entity.Conversation(conversation)
	}
	return conversations, nil
}

// Enable makes the model appear as though it is active/focussed.
func (m *ConversationsModel) Enable() {
	m.switchStyles(enabledConversationsStyles())
}

// Disable makes the model appear as though it is not active/focussed.
func (m *ConversationsModel) Disable() {
	m.switchStyles(disabledConversationsStyles())
}

// SetSize calculates and applies the correct size to its nested components.
// The given width and height should be the dimensions for this component not the window.
func (m *ConversationsModel) SetSize(width, height int) {
	m.list.SetSize(width, height)
}

// switchStyles updates the model state to have the given styles.
// It also updates the styles of nested components using the provided styles.
func (m *ConversationsModel) switchStyles(styles *conversationsStyles) {
	m.styles = styles
	m.list.Styles = styles.List
	m.list.SetDelegate(NewListDelegate(DefaultListDelegateKeyMap(), styles.ListItem))
}

func (m *ConversationsModel) View() string {
	return m.list.View()
}
