package models

import (
	"forum/pkg/validator"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword string
	Created        time.Time
}

type UserLoginForm struct {
	Email               string
	Password            string
	validator.Validator `form:"-"`
}

type UserSignupForm struct {
	Name                string
	Email               string
	Password            string
	PasswordAgain       string
	validator.Validator `form:"-"`
}

type UserReaction struct {
	Reaction        string
	CommentReaction string
	Comment         string
	CommentID       int
}
