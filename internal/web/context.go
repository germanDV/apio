package web

import (
	"context"
	"net/http"
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
func GetUser(ctx context.Context) (CtxUser, *ApiErr) {
	u, ok := ctx.Value(ctxUserKey).(CtxUser)
	if !ok {
		return CtxUser{}, &ApiErr{Code: http.StatusUnauthorized, Msg: "no user in context"}
	}
	return u, nil
}
