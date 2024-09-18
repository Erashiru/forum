package service

import (
	"forum/internal/repo"
	"forum/models"
	"net/http"
)

type service struct {
	repo repo.RepoI
}

type UserService interface {
	UserSignUp(data *models.TemplateData, form *models.UserSignupForm) (*models.TemplateData, error)
	UserLogin(data *models.TemplateData, form *models.UserLoginForm) (*models.TemplateData, int, error)
	SaveUser(form *models.UserSignupForm) (int, error) // Add this
	GetUserByEmail(email string) (*models.User, error) // Add this
	GetUserData(data *models.TemplateData) (*models.TemplateData, error)
}

type HomePage interface {
	MainLoader(data *models.TemplateData) (*models.TemplateData, error)
	HomeUpdates(data *models.TemplateData, filterData *models.FilterData, postID string, isStandardPath bool) (*models.TemplateData, error)
	LikedPostsLoader(data *models.TemplateData) (*models.TemplateData, error)
	LikedPostsUpdate(data *models.TemplateData, postID int, reaction string) (*models.TemplateData, error)
	DislikedPostsLoader(data *models.TemplateData) (*models.TemplateData, error)
	DislikedPostsUpdate(data *models.TemplateData, postID int, reaction string) (*models.TemplateData, error)
}

type PostsService interface {
	LatestPosts() ([]*models.Posts, error)
	ReactionDone(used_id, post_id, reaction int) error
	GetPostData(data *models.TemplateData, postId int) (*models.TemplateData, error)
	PostUpdate(data *models.TemplateData, userReactions *models.UserReaction, postID int) (*models.TemplateData, error)
	PostCreate(data *models.TemplateData, Form *models.PostForm) (*models.TemplateData, int, error)
	UpdatePost(data *models.TemplateData, form *models.PostForm) error
	DeletePost(postID int) error
}

type Moderation interface {
	ModerationRequest(id int) error
	CreateReport(report *models.Report) error
}

type SessionManager interface {
	IsValidToken(token string) (bool, error)
	GetUserIDBySessionToken(sessionToken string) int
	CreateSession(w http.ResponseWriter, r *http.Request, UserID int) error
	DeleteSession(sessionID string) error
}

type UserPosts interface {
	UserPosts(data *models.TemplateData) (*models.TemplateData, error)
	UserPostsUpdate(data *models.TemplateData, postID int, reaction string) (*models.TemplateData, error)
}

type AdminService interface {
	ReportsList() ([]*models.Report, error)
	RequestsList() ([]*models.Request, error)
	UpdateUserRole(data *models.TemplateData, role string) error
	ReportDone(id int) error
	PostDelete(postid int) error
}

type ServiceI interface {
	UserPosts
	SessionManager
	PostsService
	HomePage
	UserService
	Moderation
	AdminService
}

func New(r repo.RepoI) ServiceI {
	return &service{
		r,
	}
}
