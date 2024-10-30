package views

import (
	"github.com/Broderick-Westrope/teatime/internal/tui/components"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

func defaultAppKeyMaps() map[appFocusRegion]components.KeyMap {
	toggleHelp := key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help"))
	forceQuit := key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "force quit"))

	return map[appFocusRegion]components.KeyMap{
		appFocusRegionChat: &chatKeyMap{
			Back: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
			Send: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "send")),

			ToggleHelp: toggleHelp,
			ForceQuit:  forceQuit,
		},
		appFocusRegionConversations: &conversationsKeyMap{
			Quit: key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),

			List: list.KeyMap{
				// Browsing.
				CursorUp: key.NewBinding(
					key.WithKeys("up", "k"),
					key.WithHelp("↑/k", "up"),
				),
				CursorDown: key.NewBinding(
					key.WithKeys("down", "j"),
					key.WithHelp("↓/j", "down"),
				),
				PrevPage: key.NewBinding(
					key.WithKeys("left", "h", "pgup", "b", "u"),
					key.WithHelp("←/h/pgup", "prev page"),
				),
				NextPage: key.NewBinding(
					key.WithKeys("right", "l", "pgdown", "f", "d"),
					key.WithHelp("→/l/pgdn", "next page"),
				),
				GoToStart: key.NewBinding(
					key.WithKeys("home", "g"),
					key.WithHelp("g/home", "go to start"),
				),
				GoToEnd: key.NewBinding(
					key.WithKeys("end", "G"),
					key.WithHelp("G/end", "go to end"),
				),
				Filter: key.NewBinding(
					key.WithKeys("/"),
					key.WithHelp("/", "filter"),
				),
				ClearFilter: key.NewBinding(
					key.WithKeys("esc"),
					key.WithHelp("esc", "clear filter"),
				),

				// Filtering.
				CancelWhileFiltering: key.NewBinding(
					key.WithKeys("esc"),
					key.WithHelp("esc", "cancel"),
				),
				AcceptWhileFiltering: key.NewBinding(
					key.WithKeys("enter", "tab", "shift+tab", "ctrl+k", "up", "ctrl+j", "down"),
					key.WithHelp("enter", "apply filter"),
				),

				// Toggle help.
				ShowFullHelp: key.NewBinding(
					key.WithKeys("?"),
					key.WithHelp("?", "more"),
				),
				CloseFullHelp: key.NewBinding(
					key.WithKeys("?"),
					key.WithHelp("?", "close help"),
				),

				// Quitting. Handled by the starter model.
				Quit:      key.NewBinding(key.WithDisabled()),
				ForceQuit: key.NewBinding(key.WithDisabled()),
			},
			ListDelegate: &components.ConversationListDelegateKeyMap{
				Open: key.NewBinding(
					key.WithKeys("enter"),
					key.WithHelp("enter", "open"),
				),
				New: key.NewBinding(
					key.WithKeys("n", "+"),
					key.WithHelp("n/+", "new"),
				),
				Delete: key.NewBinding(
					key.WithKeys("backspace"),
					key.WithHelp("delete/backspace", "delete"),
				),
			},

			ToggleHelp: toggleHelp,
			ForceQuit:  forceQuit,
		},
	}
}

// Chat -----------------------------------------------------------------------

type chatKeyMap struct {
	// Contextual
	Back key.Binding // Remove focus from the chat component.
	Send key.Binding // Send the written message.

	// Persistent
	ToggleHelp key.Binding
	ForceQuit  key.Binding
}

func (k *chatKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Back,
		k.Send,
	}
}

func (k *chatKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Back,
			k.Send,
		},
		{
			k.ToggleHelp,
			k.ForceQuit,
		},
	}
}

// Conversations --------------------------------------------------------------

type conversationsKeyMap struct {
	// Contextual
	Quit key.Binding // Quit the application.

	// Dependencies
	List         list.KeyMap
	ListDelegate *components.ConversationListDelegateKeyMap

	// Persistent
	ToggleHelp key.Binding
	ForceQuit  key.Binding
}

func (k *conversationsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.List.CursorUp,
		k.List.CursorDown,

		k.ListDelegate.New,
		k.ListDelegate.Open,
		k.ListDelegate.Delete,

		k.List.Filter,
		k.List.ClearFilter,
		k.List.AcceptWhileFiltering,
		k.List.CancelWhileFiltering,

		k.ToggleHelp,
		k.Quit,
	}
}

func (k *conversationsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.List.CursorUp,
			k.List.CursorDown,
			k.List.NextPage,
			k.List.PrevPage,
			k.List.GoToStart,
			k.List.GoToEnd,
		},
		{
			k.ListDelegate.New,
			k.ListDelegate.Open,
			k.ListDelegate.Delete,
			k.Quit,
		},
		{
			k.List.Filter,
			k.List.ClearFilter,
			k.List.AcceptWhileFiltering,
			k.List.CancelWhileFiltering,
		},
		{
			k.ToggleHelp,
			k.ForceQuit,
		},
	}
}
