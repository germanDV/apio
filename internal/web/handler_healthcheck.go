package web

import (
	"net/http"
	"runtime"
	"runtime/debug"

	"github.com/a-h/rest"
)

type healthcheckResp struct {
	Status    string `json:"status"`
	GoVersion string `json:"go_version"`
	Revision  string `json:"revision"`
}

func (api *Api) handleHealthcheck() func(http.ResponseWriter, *http.Request) error {
	api.oas.Get("/healthcheck").
		HasResponseModel(http.StatusOK, rest.ModelOf[healthcheckResp]()).
		HasTags([]string{"Observability"}).
		HasOperationID("getHealthcheck").
		HasDescription("Get a minimal healthcheck response.")

	return func(w http.ResponseWriter, _ *http.Request) error {
		goVersion := runtime.Version()

		revision := "unknown"
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					revision = setting.Value
					break
				}
			}
		}

		WriteJSON(w, healthcheckResp{Status: "ok", GoVersion: goVersion, Revision: revision}, http.StatusOK)
		return nil
	}
}
