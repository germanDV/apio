package web

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/germandv/apio/internal/config"
	"github.com/germandv/apio/internal/tokenauth"
)

// Api is the main API HTTP server wrapper.
type Api struct {
	mux    *http.ServeMux
	server *http.Server
	logger *slog.Logger
}

// New creates a new Api instance.
func New(logger *slog.Logger, auth tokenauth.Service) *Api {
	cfg := config.Get()
	mux := &http.ServeMux{}

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           middleware(mux, logger, auth),
		IdleTimeout:       1 * time.Minute,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	return &Api{
		mux:    mux,
		server: server,
		logger: logger,
	}
}

// ListenAndServe starts the server in a goroutine.
// It blocks until the server is shut down and it handles graceful shutdown.
func (api *Api) ListenAndServe() {
	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		err := api.server.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				api.logger.Info("Server closed")
			} else {
				panic(err)
			}
		}
	}()

	api.logger.Info("Server running", "addr", api.server.Addr)
	api.logger.Info("Open API Spec available", "route", "/swagger-ui/swagger.json")
	api.logger.Info("Open API UI available", "route", "/swagger-ui")

	<-killSig
	api.logger.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := api.server.Shutdown(ctx)
	if err != nil {
		panic(err)
	}
}

// Route adds a new route to the API, with its handler and optional middleware
func (api *Api) Route(
	path string,
	handler func(w http.ResponseWriter, r *http.Request) error,
	middleware ...func(http.Handler) http.Handler,
) {
	var h http.Handler = ApiFunc{handler: handler, logger: api.logger}
	for _, m := range middleware {
		h = m(h)
	}
	api.mux.Handle(path, h)
}
