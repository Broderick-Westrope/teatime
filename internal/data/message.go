package data

import "time"

type Message struct {
	Content string
	Author  string
	SentAt  time.Time
}
