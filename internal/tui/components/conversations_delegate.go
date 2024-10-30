package components

import (
	"fmt"

	"github.com/Broderick-Westrope/teatime/internal/data"
	"github.com/Broderick-Westrope/teatime/internal/tui"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Conversation data.Conversation

func (c Conversation) Title() string       { return c.Name }
func (c Conversation) Description() string { return c.Messages[len(c.Messages)-1].Content }
func (c Conversation) FilterValue() string { return c.Name }

func NewListDelegate(keys *ConversationListDelegateKeyMap, styles list.DefaultItemStyles) list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles = styles

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		contact, ok := m.SelectedItem().(Conversation)
		if !ok {
			return tui.FatalErrorCmd(fmt.Errorf(
				"list delegate failed to get selected item: %w",
				tui.ErrInvalidTypeAssertion,
			))
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.Open):
				return tui.SetConversationCmd(data.Conversation(contact))

			case key.Matches(msg, keys.New):
				return m.NewStatusMessage("Creating new")

			case key.Matches(msg, keys.Delete):
				if m.FilterState() == list.FilterApplied {
					return nil
				}
				m.RemoveItem(m.Index())
				if len(m.Items()) == 0 {
					keys.Delete.SetEnabled(false)
				}
				return m.NewStatusMessage("Deleted " + contact.Name)
			}
		}
		return nil
	}
	return d
}

type ConversationListDelegateKeyMap struct {
	Open   key.Binding
	New    key.Binding
	Delete key.Binding
}
