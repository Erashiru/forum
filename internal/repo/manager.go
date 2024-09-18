package repo

import (
	"forum/internal/repo/database"
	"forum/models"
	"net/http"
)

type PostsRepo interface {
	CreatePost(title string, content string, userid int) (int, error)
	GetPost(id int) (*models.Posts, error)
	UserPosts(userid int) ([]*models.Posts, error)
	LatestPosts() ([]*models.Posts, error)
	GetLikedPost(userid int) ([]*models.Posts, error)
	GetDislikedPost(userid int) ([]*models.Posts, error)
	DeleteUserPost(postid int) error
	UpdateUserPost(postid int, text string, title string) error
}

type UsersRepo interface {
	CreateUser(username, email, password string) (int, error)
	Authenticate(email, password string) (int, error)
	Exitsts(name string) (bool, error)
	GetUser(id int) (string, error)
	SaveUser(form *models.UserSignupForm) (int, error)
	GetUserByEmail(email string) (*models.User, error)
	GetRole(id int) (string, error)
	CreateUserRole(id int, role string) error
}

type ReactionsRepo interface {
	CreateReaction(userid, postid, reaction int) error
	GetLikes(postid int) (int, error)
	GetDislikes(postid int) (int, error)
}

type CommentsRepo interface {
	CreateComment(postid, userid int, text string) (int, error)
	GetComment(id int) (*models.Comment, error)
	GetComments(id int) ([]*models.Comment, error)
}

type SessionsRepo interface {
	IsValidToken(token string) (bool, error)
	GetUserIDBySessionToken(sessionToken string) int
	DeleteSession(sessionID string) error
	CreateSession(w http.ResponseWriter, r *http.Request, UserID int) error
}

type CommentReactRepo interface {
	CreateCommentReaction(userid, commentid, reaction int) error
	GetCommentLikes(commentid int) (int, error)
	GetCommentDislikes(commentid int) (int, error)
}

type Categories interface {
	ChooseCategories(postid int, categorie []string) error
	GetCategory(postid int) ([]string, error)
}

type ModerationRepo interface {
	CreateRequest(userid int) error
	RequestDone(id int) error
	CreateReport(report *models.Report) error
}

type AdminRepo interface {
	UpdateUserRole(id int, role string) error
	GetRequests() ([]*models.Request, error)
	GetReports() ([]*models.Report, error)
	ReportDone(id int) error
	PostDelete(postID int) error
}

type RepoI interface {
	PostsRepo
	UsersRepo
	ReactionsRepo
	CommentsRepo
	SessionsRepo
	CommentReactRepo
	Categories
	ModerationRepo
	AdminRepo
}

func New(storagePath string) (RepoI, error) {
	return database.New(storagePath)
}
