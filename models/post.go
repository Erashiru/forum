package models

import (
	"forum/pkg/validator"
	"time"
)

type Posts struct {
	ID       int
	Title    string
	Content  string
	Created  time.Time
	Username string
	Comments []*Comment
	Likes    int
	Dislikes int
	Category []string
	URL      string
}

type Reaction struct {
	ReactionID     int
	UserID         int
	PostID         int
	ReactionStatus int
}

type Comment struct {
	CommentId int
	UserID    int
	Username  string
	PostID    int
	Text      string
	Created   time.Time
	Likes     int
	Dislikes  int
}

type CommentForm struct {
	Text   string
	Userid int
	validator.Validator
}

type PostForm struct {
	Title               string
	Content             string
	Category            []string
	validator.Validator `form:"-"`
}

type FilterData struct {
	Category    string
	Filter      string
	Reaction    string
	CategoryURL string
}
