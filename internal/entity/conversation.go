package entity

// Conversation is a list of messages between several participants.
type Conversation struct {
	Name         string    `json:"name"`
	Participants []string  `json:"participants"`
	Messages     []Message `json:"messages"`
}
