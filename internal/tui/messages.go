package tui

import (
	"github.com/Broderick-Westrope/teatime/internal/entity"
)

// FatalErrorMsg encloses an error which should be set on the starter model before exiting the program.
type FatalErrorMsg error

// CreateConversationMsg encloses the details for creating a new conversation.
type CreateConversationMsg struct {
	Name               string
	Participants       []string
	NotifyParticipants bool
}

// SetConversationMsg encloses the contact whose conversation should be displayed the chat.
type SetConversationMsg entity.Conversation

// SendMessageMsg encloses a new message that needs to be persisted locally and sent to the conversation participants.
type SendMessageMsg struct {
	Message      entity.Message
	Conversation entity.Conversation
}

// ReceiveMessageMsg encloses a new message that needs to be persisted locally.
type ReceiveMessageMsg struct {
	ConversationName string
	Message          entity.Message
}

// OpenModalMsg encloses a modal which should be opened on top of the current content.
type OpenModalMsg struct {
	Modal Modal
}

// CloseModalMsg signals that any open modals should be closed.
type CloseModalMsg struct{}

type QuitMsg struct{}

type DebugLogMsg string
