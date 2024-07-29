package main

import (
	"log/slog"
	"net/http"

	"github.com/germandv/apio/internal/config"
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

func handleTestAuth(logger *slog.Logger) func(http.ResponseWriter, *http.Request) *web.ApiErr {
	return func(w http.ResponseWriter, r *http.Request) *web.ApiErr {
		user, err := web.GetUser(r.Context())
		if err != nil {
			return err
		}

		if user.Role != "admin" {
			return &web.ApiErr{Msg: "you have no access to this", Code: http.StatusForbidden}
		}

		env := web.Envelope{"message": "you have access to this", "user": user.ID}
		web.WriteJSON(w, env, http.StatusOK)
		return nil
	}
}

func handleHealthcheck(logger *slog.Logger) func(http.ResponseWriter, *http.Request) *web.ApiErr {
	return func(w http.ResponseWriter, _ *http.Request) *web.ApiErr {
		web.WriteJSON(w, web.Envelope{"status": "ok"}, http.StatusOK)
		return nil
	}
}
