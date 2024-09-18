package models

import "errors"

type CustomError struct {
	ErrorCode int
	ErrorMsg  string
}

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrDuplicateName      = errors.New("models: duplicate name")
	ErrNotValidPostForm   = errors.New("hanlders: not valid create post form")
	NotFound              = errors.New("Not Found 404")
	NoPermission          = errors.New("No Permission 403")
	ErrDuplicateRequest   = errors.New("Duplicated request")
)
