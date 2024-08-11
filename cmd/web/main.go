package main

import (
	"net/http"

	"github.com/a-h/rest"

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

	oas := rest.NewAPI("apio")
	api := web.New(logger, auth)
	api.Route("GET /healthcheck", handleHealthcheck(oas, logger))
	api.Route("POST /test-auth/{id}", handleTestAuth(oas, logger))

	// Handler for Swagger UI.
	oasUIHandler, err := setupOpenApiSpec(oas)
	if err != nil {
		panic(err)
	}
	api.Route("GET /swagger-ui/*", func(w http.ResponseWriter, r *http.Request) error {
		oasUIHandler.ServeHTTP(w, r)
		return nil
	})

	api.ListenAndServe()
}
