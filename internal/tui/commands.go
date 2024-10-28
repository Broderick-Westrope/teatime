package tui

import (
	"github.com/Broderick-Westrope/teatime/internal/data"
	tea "github.com/charmbracelet/bubbletea"
)

// FatalErrorCmd returns a command for creating a new FatalErrorMsg with the given error.
func FatalErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return FatalErrorMsg(err)
	}
}

// SetConversationCmd is a command for creating a new SetConversationMsg.
func SetConversationCmd(contact data.Conversation) tea.Cmd {
	return func() tea.Msg {
		return SetConversationMsg(contact)
	}
}

// SendMessageCmd is a command for creating a new SendMessageMsg.
func SendMessageCmd(message data.Message, conversation data.Conversation) tea.Cmd {
	return func() tea.Msg {
		return SendMessageMsg{
			Message:      message,
			Conversation: conversation,
		}
	}
}

func DebugLogCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return DebugLogMsg(value)
	}
}
