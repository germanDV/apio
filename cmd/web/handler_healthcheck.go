package main

import (
	"log/slog"
	"net/http"

	"github.com/a-h/rest"

	"github.com/germandv/apio/internal/web"
)

type healthcheckResp struct {
	Status string `json:"status"`
}

func handleHealthcheck(oas *rest.API, _ *slog.Logger) func(http.ResponseWriter, *http.Request) error {
	oas.Get("/healthcheck").
		HasResponseModel(http.StatusOK, rest.ModelOf[healthcheckResp]()).
		HasTags([]string{"Observability"}).
		HasOperationID("getHealthcheck").
		HasDescription("Get a minimal healthcheck response.")

	return func(w http.ResponseWriter, _ *http.Request) error {
		web.WriteJSON(w, healthcheckResp{Status: "ok"}, http.StatusOK)
		return nil
	}
}
