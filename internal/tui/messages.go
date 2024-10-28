package tui

import "github.com/Broderick-Westrope/teatime/internal/data"

// FatalErrorMsg encloses an error which should be set on the starter model before exiting the program.
type FatalErrorMsg error

// SetConversationMsg encloses the contact whose conversation should be displayed the chat.
type SetConversationMsg data.Conversation

// SendMessageMsg encloses a new message that needs to be persisted locally and sent to the conversation participants.
type SendMessageMsg struct {
	Message      data.Message
	Conversation data.Conversation
}

// ReceiveMessageMsg encloses a new message that needs to be persisted locally.
type ReceiveMessageMsg struct {
	ConversationName string
	Message          data.Message
}

type DebugLogMsg string
