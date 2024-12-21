package web

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/germandv/apio/internal/errs"
)

// APIFunc holds an http.HandlerFunc that returns an error.
type APIFunc struct {
	handler func(w http.ResponseWriter, r *http.Request) error
	logger  *slog.Logger
}

func (af APIFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	status := http.StatusOK

	// Handle panic
	defer func() {
		err := recover()
		if err != nil {
			WriteJSON(w, Envelope{"error": "internal server error"}, http.StatusInternalServerError)
			af.logger.Warn(
				"ApiFunc panicked",
				"method", r.Method,
				"path", r.URL.Path,
				"status", http.StatusInternalServerError,
				"took_ms", time.Since(start).Milliseconds(),
				"ip", r.RemoteAddr,
				"err", err,
			)
		}
	}()

	err := af.handler(w, r)
	errMsg := ""

	// Map errors to HTTP status codes.
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoUserInCtx):
			errMsg = err.Error()
			status = http.StatusUnauthorized
		case errors.Is(err, errs.ErrNoPermission):
			errMsg = err.Error()
			status = http.StatusForbidden
		case errors.Is(err, errs.ErrDuplicateTag):
			errMsg = err.Error()
			status = http.StatusConflict
		case errors.Is(err, errs.ErrInvalidID):
		case errors.Is(err, errs.ErrMaxLen):
		case errors.Is(err, errs.ErrTagNotFound):
		case errors.Is(err, errs.ErrEmptyTitle):
		case errors.Is(err, errs.ErrEmptyName):
			errMsg = err.Error()
			status = http.StatusBadRequest
		default:
			errMsg = "internal server error" // overwrite error message to avoid leaking internal details.
			status = http.StatusInternalServerError
		}

		// Send error response.
		WriteJSON(w, errs.ErrResp{Error: errMsg}, status)
	}

	af.logger.Info(
		"API handler finished",
		"method", r.Method,
		"path", r.URL.Path,
		"status", status,
		"took_ms", time.Since(start).Milliseconds(),
		"ip", r.RemoteAddr,
		"err", errMsg,
	)
}
