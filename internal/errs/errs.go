package errs

import "errors"

var (
	ErrNoUserInCtx  = errors.New("no user in context")
	ErrNoPermission = errors.New("you have no access to this resource")
)
