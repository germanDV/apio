package errs

import "errors"

type ErrResp struct {
	Error string `json:"error"`
}

var (
	ErrNoUserInCtx  = errors.New("no user in context")
	ErrNoPermission = errors.New("you have no access to this resource")
	ErrDuplicateTag = errors.New("duplicate tag")
	ErrInvalidID    = errors.New("invalid ID")
	ErrEmptyName    = errors.New("name is required")
	ErrEmptyTitle   = errors.New("title is required")
	ErrMinLen       = errors.New("min length not met")
	ErrMaxLen       = errors.New("max length exceeded")
	ErrTagNotFound  = errors.New("one or more tags do not exist")
)
