package models

import "time"

type Request struct {
	ID         int
	SenderID   int
	SenderName string
	CreatedAt  time.Time
}
