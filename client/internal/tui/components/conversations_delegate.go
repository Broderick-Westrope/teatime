package components

import (
	"fmt"

	"github.com/Broderick-Westrope/teatime/client/internal/tui"
	"github.com/Broderick-Westrope/teatime/client/internal/tui/modals"
	"github.com/Broderick-Westrope/teatime/internal/entity"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Conversation entity.Conversation

func (c Conversation) Title() string       { return c.Metadata.Name }
func (c Conversation) FilterValue() string { return c.Metadata.Name }
func (c Conversation) Description() string {
	if len(c.Messages) <= 0 {
		return ""
	}
	return c.Messages[len(c.Messages)-1].Content
}

func NewListDelegate(keys *ListDelegateKeyMap, styles list.DefaultItemStyles) list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles = styles

	d.ShortHelpFunc = keys.ShortHelp
	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{keys.ShortHelp()}
	}

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		// Operations that do not require a selected conversation
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.new):
				return tui.OpenModalCmd(modals.NewCreateConversationModel())
			}
		}

		selectedItem := m.SelectedItem()
		if selectedItem == nil {
			return nil
		}

		conversation, ok := selectedItem.(Conversation)
		if !ok {
			return tui.FatalErrorCmd(fmt.Errorf(
				"list delegate failed to get selected item: %w",
				tui.ErrInvalidTypeAssertion,
			))
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.submit):
				return tui.SetConversationCmd(entity.Conversation(conversation))

			case key.Matches(msg, keys.delete):
				// TODO: use a cmd to open a confirmation modal before deleting this conversation
				if m.FilterState() == list.FilterApplied {
					return nil
				}

				//m.RemoveItem(m.Index())
				return tui.OpenModalCmd(modals.NewDeleteConversationModel(conversation.Metadata))
			}
		}
		return nil
	}
	return d
}

type ListDelegateKeyMap struct {
	submit key.Binding
	new    key.Binding
	delete key.Binding
}

func DefaultListDelegateKeyMap() *ListDelegateKeyMap {
	return &ListDelegateKeyMap{
		submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select contact"),
		),
		new: key.NewBinding(
			key.WithKeys("n", "+"),
			key.WithHelp("n/+", "new contact"),
		),
		delete: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("delete/backspace", "delete contact"),
		),
	}
}

func (d ListDelegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.submit,
		d.new,
		d.delete,
	}
}
