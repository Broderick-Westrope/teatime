package data

import "time"

type Message struct {
	Content string    `json:"content"`
	Author  string    `json:"author"`
	SentAt  time.Time `json:"sent_at"`
}
