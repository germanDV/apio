package web

import (
	"log/slog"
	"net/http"
	"time"
)

// ApiFunc is an http.HandlerFunc that returns an ApiErr
type ApiFunc struct {
	handler func(w http.ResponseWriter, r *http.Request) *ApiErr
	logger  *slog.Logger
}

func (af ApiFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	if err != nil {
		msg := err.Msg
		status = err.Code

		// Overwrite error message to avoid leaking internal details.
		if status >= 500 {
			msg = "internal server error"
		}

		// Send error response
		WriteJSON(w, Envelope{"error": msg}, status)
		errMsg = err.Msg
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
