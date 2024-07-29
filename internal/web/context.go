package web

import (
	"context"

	"github.com/germandv/apio/internal/errs"
)

// CtxUser is the user in the request context,
// created from the claims in the JWT.
type CtxUser struct {
	ID   string
	Role string
}

type CtxKey string

const ctxUserKey CtxKey = "ctx-user-key"

// GetUser returns the user from the request context.
func GetUser(ctx context.Context) (CtxUser, error) {
	u, ok := ctx.Value(ctxUserKey).(CtxUser)
	if !ok {
		return CtxUser{}, errs.ErrNoUserInCtx
	}
	return u, nil
}
