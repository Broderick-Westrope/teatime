package data

import "time"

// Message is a single chat message sent from a user.
type Message struct {
	Content string    `json:"content"`
	Author  string    `json:"author"`
	SentAt  time.Time `json:"sent_at"`
}
