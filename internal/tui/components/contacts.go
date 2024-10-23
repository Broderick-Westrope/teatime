package components

import (
	"fmt"

	"github.com/Broderick-Westrope/teatime/internal/data"
	"github.com/Broderick-Westrope/teatime/internal/tui"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = &ContactsModel{}

type ContactsModel struct {
	list   list.Model
	styles *contactsStyles
}

func NewContactsModel(contacts []data.Contact, enabled bool) *ContactsModel {
	var items = make([]list.Item, len(contacts))
	for i, d := range contacts {
		items[i] = Contact(d)
	}

	styles := disabledContactsStyles()
	if enabled {
		styles = enabledContactsStyles()
	}

	delegate := NewListDelegate(DefaultListDelegateKeyMap(), styles.ListItem)
	contactList := list.New(items, delegate, 0, 0)
	contactList.Title = "Contacts"

	return &ContactsModel{
		list:   contactList,
		styles: styles,
	}
}

func (m *ContactsModel) Init() tea.Cmd {
	m.switchStyles(m.styles)
	return nil
}

func (m *ContactsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *ContactsModel) View() string {
	return m.list.View()
}

func (m *ContactsModel) GetSelectedContact() (*Contact, error) {
	contact, ok := m.list.SelectedItem().(Contact)
	if !ok {
		return nil, fmt.Errorf("failed to get selected contact: %w", tui.ErrInvalidTypeAssertion)
	}
	return &contact, nil
}

func (m *ContactsModel) AddNewMessage(in tui.SendMessageMsg) (tea.Cmd, error) {
	foundIdx := -1
	items := m.list.Items()
	for i, item := range items {
		contact, ok := item.(Contact)
		if !ok {
			return nil, fmt.Errorf("failed to add new message to contacts: %v", tui.ErrInvalidTypeAssertion)
		}

		if contact.Username != in.RecipientUsername {
			continue
		}

		foundIdx = i
		contact.Conversation = append(contact.Conversation, in.Message)
		items[i] = contact
		break
	}

	switch {
	case foundIdx < 0:
		return nil, fmt.Errorf("failed to add new message to contacts: could not find message recipient")
	case foundIdx != 0:
		item := items[foundIdx]
		items = append(items[:foundIdx], items[foundIdx+1:]...)
		items = append([]list.Item{item}, items...)
		m.list.Select(0)
	}
	return m.list.SetItems(items), nil
}

func (m *ContactsModel) Enable() {
	m.switchStyles(enabledContactsStyles())
}

func (m *ContactsModel) Disable() {
	m.switchStyles(disabledContactsStyles())
}

func (m *ContactsModel) switchStyles(styles *contactsStyles) {
	m.styles = styles
	m.list.Styles = styles.List
	m.list.SetDelegate(NewListDelegate(DefaultListDelegateKeyMap(), styles.ListItem))
}

func (m *ContactsModel) SetSize(width, height int) {
	m.list.SetSize(width, height)
}
