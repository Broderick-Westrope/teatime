package tui

// ComponentSizeMsg encloses dimensions that should be interpreted by the child component sensibly.
// This is used when a top-level view receives a window resize message and wants to resize a child component.
type ComponentSizeMsg struct {
	Width  int
	Height int
}

// FatalErrorMsg encloses an error which should be set on the starter model before exiting the program.
type FatalErrorMsg error

// UpdateChatMsg signifies a need to update the chat to use the conversation of the currently selected contact.
type UpdateChatMsg struct{}

type DebugLogMsg string
