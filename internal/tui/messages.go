package tui

import "github.com/Broderick-Westrope/teatime/internal/entity"

// FatalErrorMsg encloses an error which should be set on the starter model before exiting the program.
type FatalErrorMsg error

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

type QuitMsg struct{}

type DebugLogMsg string
