package models

import "time"

type Report struct {
	PostID    int
	Text      string
	ModerID   int
	ModerName string
	CreatedAt time.Time
}
