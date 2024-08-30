package web

import (
	"net/http"

	"github.com/a-h/rest"
	"github.com/a-h/rest/swaggerui"
	"github.com/germandv/apio/internal/errs"
	"github.com/getkin/kin-openapi/openapi3"
)

func (api *Api) setupOpenApiSpec() (http.Handler, error) {
	api.oas.StripPkgPaths = []string{"main", "github.com/germandv/apio/internal"}

	// Common error response to all routes.
	for _, rr := range api.oas.Routes {
		for _, op := range rr {
			op.HasResponseModel(http.StatusUnauthorized, rest.ModelOf[errs.ErrResp]())
			op.HasResponseModel(http.StatusTooManyRequests, rest.ModelOf[errs.ErrResp]())
			op.HasResponseModel(http.StatusInternalServerError, rest.ModelOf[errs.ErrResp]())
		}
	}

	spec, err := api.oas.Spec()
	if err != nil {
		return nil, err
	}

	spec.Info.Title = "Apio"
	spec.Info.Version = "v1.0.0"
	spec.Info.Description = "Documentation for APIO"

	// Configure ability to provide Authorization header.
	spec.Components.SecuritySchemes = map[string]*openapi3.SecuritySchemeRef{
		"bearerAuth": {
			Value: openapi3.NewJWTSecurityScheme(),
		},
	}
	securitySchemeToScopes := openapi3.NewSecurityRequirement()
	securitySchemeToScopes.Authenticate("bearerAuth")
	spec.Security = *openapi3.NewSecurityRequirements().With(securitySchemeToScopes)

	// Generate handler for Swgger UI.
	uiHandler, err := swaggerui.New(spec)
	if err != nil {
		return nil, err
	}

	return uiHandler, nil
}
