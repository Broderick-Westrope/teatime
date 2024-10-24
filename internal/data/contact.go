package data

// Contact is another user and this users' conversation with them.
type Contact struct {
	Username     string    `json:"username"`
	Conversation []Message `json:"conversation"`
}
