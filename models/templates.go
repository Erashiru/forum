package models

type TemplateData struct {
	UserID      int
	CurrentYear int
	Post        Posts
	Posts       []*Posts
	Form        any
	User        User
	Username    string
	Comments    []*Comment
	Likes       int
	Dislikes    int
	Requests    []*Request
	Role        string
	ReqUsrID    int
	Reports     []*Report
	PostID      int
}
