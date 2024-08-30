package web

import "net/http"

func (api *Api) routes() {
	// API endpoints.
	api.Route("GET /healthcheck", api.handleHealthcheck())
	api.Route("POST /tags", api.handleCreateTag(), RequireUser)
	api.Route("GET /tags", api.handleGetTags(), RequireUser)
	api.Route("POST /notes", api.handleCreateNote(), RequireUser)
	api.Route("GET /notes", api.handleGetNotes(), RequireUser)

	// Handlers for Swagger UI.
	oasUIHandler, err := api.setupOpenApiSpec()
	if err != nil {
		panic(err)
	}
	api.Route("GET /swagger-ui", func(w http.ResponseWriter, r *http.Request) error {
		oasUIHandler.ServeHTTP(w, r)
		return nil
	})
	api.Route("GET /swagger-ui/", func(w http.ResponseWriter, r *http.Request) error {
		oasUIHandler.ServeHTTP(w, r)
		return nil
	})
}
