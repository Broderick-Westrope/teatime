package tui

import "github.com/Broderick-Westrope/teatime/internal/data"

// ComponentSizeMsg encloses dimensions that should be interpreted by the child component sensibly.
// This is used when a top-level view receives a window resize message and wants to resize a child component.
type ComponentSizeMsg struct {
	Width  int
	Height int
}

// FatalErrorMsg encloses an error which should be set on the starter model before exiting the program.
type FatalErrorMsg error

// SetConversationMsg signifies a need to update the chat to use the conversation of the currently selected contact.
type SetConversationMsg struct{}

// SendMessageMsg encloses a new message that needs to be persisted locally and sent to the recipient.
type SendMessageMsg struct {
	RecipientUsername string
	Message           data.Message
}

type DebugLogMsg string
