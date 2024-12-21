package web

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-h/rest"
)

func TestHealthcheck(t *testing.T) {
	oas := rest.NewAPI("apio_test")
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	api := &API{
		oas:    oas,
		logger: logger,
	}

	h := api.handleHealthcheck()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	err := h(w, r)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("expected %q, got %q", "application/json", w.Header().Get("Content-Type"))
	}

	resp := healthcheckResp{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if resp.Status != "ok" {
		t.Fatalf("expected resp.Status ok, got %s", resp.Status)
	}
}
