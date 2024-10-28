package components

import (
	"fmt"

	"github.com/Broderick-Westrope/teatime/internal/data"
	"github.com/Broderick-Westrope/teatime/internal/tui"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = &ConversationsModel{}

type ConversationsModel struct {
	list   list.Model
	styles *conversationsStyles
}

func NewConversationsModel(conversations []data.Conversation, enabled bool) *ConversationsModel {
	var items = make([]list.Item, len(conversations))
	for i, d := range conversations {
		items[i] = Conversation(d)
	}

	styles := disabledConversationsStyles()
	if enabled {
		styles = enabledConversationsStyles()
	}

	delegate := NewListDelegate(DefaultListDelegateKeyMap(), styles.ListItem)
	contactList := list.New(items, delegate, 0, 0)
	contactList.Title = "Contacts"

	return &ConversationsModel{
		list:   contactList,
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

// AddNewMessage will add the given message to the chat with the given chatName.
// It will also move this messages to the top of the contacts list and update the list selection.
func (m *ConversationsModel) AddNewMessage(chatName string, message data.Message) (tea.Cmd, error) {
	const methodErr = "failed to add new message to contacts"
	foundIdx := -1
	items := m.list.Items()

	for i, item := range items {
		contact, ok := item.(Conversation)
		if !ok {
			return nil, fmt.Errorf("%s: %v", methodErr, tui.ErrInvalidTypeAssertion)
		}

		if contact.Name != chatName {
			continue
		}

		foundIdx = i
		contact.Messages = append(contact.Messages, message)
		items[i] = contact
		break
	}

	switch {
	case foundIdx < 0:
		return nil, fmt.Errorf("%s: could not find messages %q", methodErr, chatName)
	case foundIdx != 0: // if the contact is not already first in the list, move it to first
		item := items[foundIdx]
		items = append(items[:foundIdx], items[foundIdx+1:]...)
		items = append([]list.Item{item}, items...)
		m.list.Select(0)
	}
	return m.list.SetItems(items), nil
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
