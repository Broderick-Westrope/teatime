package contacts

import (
	"github.com/Broderick-Westrope/teatime/internal/data"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Contact data.Contact

func (c Contact) Title() string       { return c.Username }
func (c Contact) Description() string { return c.Conversation[len(c.Conversation)-1].Content }
func (c Contact) FilterValue() string { return c.Username }

func NewListDelegate(keys *ListDelegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var selectedName string

		if i, ok := m.SelectedItem().(Contact); ok {
			selectedName = i.Username
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.submit):
				return m.NewStatusMessage("You chose " + selectedName)

			case key.Matches(msg, keys.new):
				return m.NewStatusMessage("Creating new")

			case key.Matches(msg, keys.delete):
				m.RemoveItem(m.Index())
				if len(m.Items()) == 0 {
					keys.delete.SetEnabled(false)
				}
				return m.NewStatusMessage("Deleted " + selectedName)
			}
		}

		return nil
	}

	d.ShortHelpFunc = keys.ShortHelp
	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{keys.ShortHelp()}
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
