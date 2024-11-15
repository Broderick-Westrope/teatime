package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Broderick-Westrope/teatime/internal/entity"
)

// FatalErrorMsg encloses an error which should be set on the starter model before exiting the program.
type FatalErrorMsg error

// FatalErrorCmd returns a command for creating a new FatalErrorMsg with the given error.
func FatalErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return FatalErrorMsg(err)
	}
}

// AuthenticateMsg encloses the details for attempting authentication.
type AuthenticateMsg struct {
	IsSignup    bool
	Credentials *entity.Credentials
}

// AuthenticateCmd returns a command for creating a new AuthenticateMsg.
func AuthenticateCmd(isSignup bool, username, password string) tea.Cmd {
	return func() tea.Msg {
		return AuthenticateMsg{
			IsSignup: isSignup,
			Credentials: &entity.Credentials{
				Username: username,
				Password: password,
			},
		}
	}
}

// CreateConversationMsg encloses the details for creating a new conversation.
type CreateConversationMsg struct {
	Name               string
	Participants       []string
	NotifyParticipants bool
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

// DeleteConversationMsg encloses the conversation to be deleted.
type DeleteConversationMsg struct {
	ConversationMD entity.ConversationMetadata
}

// DeleteConversationCmd returns a command for creating a new CreateConversationMsg.
func DeleteConversationCmd(conversationMD entity.ConversationMetadata) tea.Cmd {
	return func() tea.Msg {
		return DeleteConversationMsg{
			ConversationMD: conversationMD,
		}
	}
}

// SetConversationMsg encloses the contact whose conversation should be displayed the chat.
type SetConversationMsg entity.Conversation

// SetConversationCmd returns a command for creating a new SetConversationMsg.
func SetConversationCmd(contact entity.Conversation) tea.Cmd {
	return func() tea.Msg {
		return SetConversationMsg(contact)
	}
}

// SendMessageMsg encloses a new message that needs to be persisted locally and sent to the conversation participants.
type SendMessageMsg struct {
	Message        entity.Message
	ConversationMD entity.ConversationMetadata
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

// ReceiveMessageMsg encloses a new message that needs to be persisted locally.
type ReceiveMessageMsg struct {
	ConversationMD entity.ConversationMetadata
	Message        entity.Message
}

// OpenModalMsg encloses a modal which should be opened on top of the current content.
type OpenModalMsg struct {
	Modal Modal
}

// OpenModalCmd returns a command for creating a new OpenModalMsg.
func OpenModalCmd(modal Modal) tea.Cmd {
	return func() tea.Msg {
		return OpenModalMsg{
			Modal: modal,
		}
	}
}

// CloseModalMsg signals that any open modals should be closed.
type CloseModalMsg struct{}

// CloseModalCmd is a command for creating a new CloseModalMsg.
func CloseModalCmd() tea.Msg {
	return CloseModalMsg{}
}

type QuitMsg struct{}

func QuitCmd() tea.Msg {
	return QuitMsg{}
}

type DebugLogMsg string

func DebugLogCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return DebugLogMsg(value)
	}
}
