package entity

import "github.com/google/uuid"

// Conversation is a list of messages between several participants.
type Conversation struct {
	Metadata ConversationMetadata
	Messages []Message `json:"messages"`
}

type ConversationMetadata struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Participants []string  `json:"participants"`
}
