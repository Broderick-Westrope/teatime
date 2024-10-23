package data

type Contact struct {
	Username     string    `json:"username"`
	Conversation []Message `json:"conversation"`
}
