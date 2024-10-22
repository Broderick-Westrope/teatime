package contacts

import (
	"github.com/Broderick-Westrope/teatime/internal/tui"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	list list.Model
}

func NewModel(contacts []Contact) *Model {
	var items = make([]list.Item, len(contacts))
	for i, d := range contacts {
		items[i] = d
	}

	delegate := NewListDelegate(DefaultListDelegateKeyMap())
	contactList := list.New(items, delegate, 0, 0)
	contactList.Title = "Contacts"

	return &Model{
		list: contactList,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *Model) View() string {
	return m.list.View()
}
