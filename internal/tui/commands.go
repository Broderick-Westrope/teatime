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

// AuthenticateCmd returns a command for creating a new AuthenticateMsg.
func AuthenticateCmd(isSignup bool, username, password string) tea.Cmd {
	return func() tea.Msg {
		return AuthenticateMsg{
			IsSignup: isSignup,
			Username: username,
			Password: password,
		}
	}
}

// CreateConversationCmd returns a command for creating a new CreateConversationMsg.
func CreateConversationCmd(name string, participants []string, notifyParticipants bool) tea.Cmd {
	return func() tea.Msg {
		return CreateConversationMsg{
			Name:               name,
			Participants:       participants,
			NotifyParticipants: notifyParticipants,
		}
	}
}

// DeleteConversationCmd returns a command for creating a new CreateConversationMsg.
func DeleteConversationCmd(conversationMD entity.ConversationMetadata) tea.Cmd {
	return func() tea.Msg {
		return DeleteConversationMsg{
			ConversationMD: conversationMD,
		}
	}
}

// SetConversationCmd returns a command for creating a new SetConversationMsg.
func SetConversationCmd(contact entity.Conversation) tea.Cmd {
	return func() tea.Msg {
		return SetConversationMsg(contact)
	}
}

// SendMessageCmd returns a command for creating a new SendMessageMsg.
func SendMessageCmd(message entity.Message, conversationMD entity.ConversationMetadata) tea.Cmd {
	return func() tea.Msg {
		return SendMessageMsg{
			Message:        message,
			ConversationMD: conversationMD,
		}
	}
}

// OpenModalCmd returns a command for creating a new OpenModalMsg.
func OpenModalCmd(modal Modal) tea.Cmd {
	return func() tea.Msg {
		return OpenModalMsg{
			Modal: modal,
		}
	}
}

// CloseModalCmd is a command for creating a new CloseModalMsg.
func CloseModalCmd() tea.Msg {
	return CloseModalMsg{}
}

func QuitCmd() tea.Msg {
	return QuitMsg{}
}

func DebugLogCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return DebugLogMsg(value)
	}
}
