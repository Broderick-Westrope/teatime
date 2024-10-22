package components

import (
	"fmt"

	"github.com/Broderick-Westrope/teatime/internal/tui"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = &ContactsModel{}

type ContactsModel struct {
	list   list.Model
	Styles *ContactsStyles
}

func NewContactsModel(contacts []Contact, enabled bool) *ContactsModel {
	var items = make([]list.Item, len(contacts))
	for i, d := range contacts {
		items[i] = d
	}

	styles := DisabledContactsStyles()
	if enabled {
		styles = EnabledContactsStyles()
	}

	delegate := NewListDelegate(DefaultListDelegateKeyMap(), styles.ListItem)
	contactList := list.New(items, delegate, 0, 0)
	contactList.Title = "Contacts"

	m := &ContactsModel{
		list: contactList,
	}
	m.SwitchStyles(styles)

	return m
}

func (m *ContactsModel) Init() tea.Cmd {
	return nil
}

func (m *ContactsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tui.ComponentSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
		return m, nil

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
	found := false
	items := m.list.Items()
	for i, item := range items {
		contact, ok := item.(Contact)
		if !ok {
			return nil, fmt.Errorf("failed to add new message to contacts: %v", tui.ErrInvalidTypeAssertion)
		}

		if contact.Username != in.RecipientUsername {
			continue
		}

		found = true
		contact.Conversation = append(contact.Conversation, in.Message)
		items[i] = contact
	}

	if !found {
		return nil, fmt.Errorf("failed to add new message to contacts: could not find message recipient")
	}

	return m.list.SetItems(items), nil
}

func (m *ContactsModel) SwitchStyles(styles *ContactsStyles) {
	m.Styles = styles
	m.list.Styles = styles.List
	m.list.SetDelegate(NewListDelegate(DefaultListDelegateKeyMap(), styles.ListItem))
}
