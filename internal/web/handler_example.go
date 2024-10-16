package web

// TODO: delete this file when examples are no longer needed.

import (
	"log/slog"
	"net/http"

	"github.com/a-h/rest"
	"github.com/germandv/apio/internal/errs"
	"github.com/getkin/kin-openapi/openapi3"
)

type TestAuthReq struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
}

type TestAuthResp struct {
	Message    string `json:"message"`
	User       string `json:"user"`
	PathValue  string `json:"path_value"`
	QueryValue string `json:"query_value,omitempty"`
	BodyFoo    string `json:"body_foo"`
	BodyBar    string `json:"body_bar"`
}

func handleTestAuth(oas *rest.API, _ *slog.Logger) func(http.ResponseWriter, *http.Request) error {
	oas.Post("/test-auth/{id}").
		HasTags([]string{"Users"}).
		HasResponseModel(http.StatusOK, rest.ModelOf[TestAuthResp]()).
		HasResponseModel(http.StatusForbidden, rest.ModelOf[errs.ErrResp]()).
		HasPathParameter("id", rest.PathParam{
			Type:        rest.PrimitiveTypeString,
			Description: "Some ID",
			ApplyCustomSchema: func(s *openapi3.Parameter) {
				s.Example = "xyz_789"
			},
		}).
		HasQueryParameter("n", rest.QueryParam{
			Type:        rest.PrimitiveTypeInteger,
			Description: "Some Number",
			Required:    false,
			ApplyCustomSchema: func(s *openapi3.Parameter) {
				s.Example = 42
			},
		}).
		HasRequestModel(rest.ModelOf[TestAuthReq]())

	return func(w http.ResponseWriter, r *http.Request) error {
		user, err := GetUser(r.Context())
		if err != nil {
			return err
		}

		if user.Role != "admin" {
			return errs.ErrNoPermission
		}

		id := r.PathValue("id")
		q := r.URL.Query()

		body := TestAuthReq{}
		err = ReadJSON(w, r, &body)
		if err != nil {
			return err
		}

		resp := TestAuthResp{
			Message:    "you have access to this",
			User:       user.ID,
			PathValue:  id,
			QueryValue: q.Get("n"),
			BodyFoo:    body.Foo,
			BodyBar:    body.Bar,
		}

		WriteJSON(w, resp, http.StatusOK)
		return nil
	}
}
