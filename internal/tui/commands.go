package tui

import (
	"github.com/Broderick-Westrope/teatime/internal/entity"
	tea "github.com/charmbracelet/bubbletea"
)

// FatalErrorCmd returns a command for creating a new FatalErrorMsg with the given error.
func FatalErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return FatalErrorMsg(err)
	}
}

// SetConversationCmd returns a command for creating a new SetConversationMsg.
func SetConversationCmd(contact entity.Conversation) tea.Cmd {
	return func() tea.Msg {
		return SetConversationMsg(contact)
	}
}

// SendMessageCmd returns a command for creating a new SendMessageMsg.
func SendMessageCmd(message entity.Message, conversation entity.Conversation) tea.Cmd {
	return func() tea.Msg {
		return SendMessageMsg{
			Message:      message,
			Conversation: conversation,
		}
	}
}

// OpenModalCmd returns a command for creating a new OpenModalMsg.
func OpenModalCmd(modal tea.Model) tea.Cmd {
	return func() tea.Msg {
		return OpenModalMsg{
			Modal: modal,
		}
	}
}

func QuitCmd() tea.Msg {
	return QuitMsg{}
}

func DebugLogCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return DebugLogMsg(value)
	}
}
