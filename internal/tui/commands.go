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
func SetConversationCmd() tea.Msg {
	return SetConversationMsg{}
}

// SendMessageCmd is a command for creating a new SendMessageMsg.
func SendMessageCmd(recipientUsername string, message data.Message) tea.Cmd {
	return func() tea.Msg {
		return SendMessageMsg{
			RecipientUsername: recipientUsername,
			Message:           message,
		}
	}
}

func DebugLogCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return DebugLogMsg(value)
	}
}
