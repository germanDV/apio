package main

import (
	"log/slog"
	"net/http"

	"github.com/germandv/apio/internal/web"
)

func handleHealthcheck(logger *slog.Logger) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, _ *http.Request) error {
		web.WriteJSON(w, web.Envelope{"status": "ok"}, http.StatusOK)
		return nil
	}
}
