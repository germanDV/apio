package errs

import "errors"

type ErrResp struct {
	Error string `json:"error"`
}

var (
	ErrNoUserInCtx  = errors.New("no user in context")
	ErrNoPermission = errors.New("you have no access to this resource")
)
