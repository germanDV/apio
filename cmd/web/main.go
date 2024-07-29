package main

import (
	"log/slog"
	"net/http"

	"github.com/germandv/apio/internal/config"
	"github.com/germandv/apio/internal/errs"
	"github.com/germandv/apio/internal/logger"
	"github.com/germandv/apio/internal/tokenauth"
	"github.com/germandv/apio/internal/web"
)

func main() {
	cfg := config.Get()

	logger, err := logger.Get(cfg.LogFormat, cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	auth, err := tokenauth.New(cfg.AuthPrivKey, cfg.AuthPublKey)
	if err != nil {
		panic(err)
	}

	api := web.New(logger, auth)
	api.Route("GET /healthcheck", handleHealthcheck(logger))
	api.Route("GET /test-auth", handleTestAuth(logger))

	api.ListenAndServe()
}

func handleTestAuth(logger *slog.Logger) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		user, err := web.GetUser(r.Context())
		if err != nil {
			return err
		}

		if user.Role != "admin" {
			return errs.ErrNoPermission
		}

		env := web.Envelope{"message": "you have access to this", "user": user.ID}
		web.WriteJSON(w, env, http.StatusOK)
		return nil
	}
}
