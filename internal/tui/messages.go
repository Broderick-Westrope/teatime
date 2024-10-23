package tui

import "github.com/Broderick-Westrope/teatime/internal/data"

// FatalErrorMsg encloses an error which should be set on the starter model before exiting the program.
type FatalErrorMsg error

// SetConversationMsg encloses the contact whose conversation should be displayed the chat.
type SetConversationMsg data.Contact

// SendMessageMsg encloses a new message that needs to be persisted locally and sent to the recipient.
type SendMessageMsg struct {
	RecipientUsername string
	Message           data.Message
}

type DebugLogMsg string
