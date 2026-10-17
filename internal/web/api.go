package web

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// Matomo is a Matomo instance and site the SPA reports page views to.
type Matomo struct {
	URL    string `json:"url" doc:"Base URL of the Matomo instance, ending in a slash"`
	SiteID int    `json:"siteId"`
}

// VersionOutput is the response of GET /api/version.
type VersionOutput struct {
	Body struct {
		Version     string  `json:"version"`
		Environment string  `json:"environment,omitempty" required:"false" doc:"Deployment name such as staging; the SPA marks everything but production"`
		Matomo      *Matomo `json:"matomo,omitempty" required:"false" doc:"Web analytics; absent when not configured"`
	}
}

// registerSystem registers the version endpoint.
func (s *Server) registerSystem(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-version",
		Method:      http.MethodGet,
		Path:        "/api/version",
		Summary:     "Server version and web analytics settings",
		Tags:        []string{"system"},
	}, func(ctx context.Context, _ *struct{}) (*VersionOutput, error) {
		out := &VersionOutput{}
		out.Body.Version = s.Version
		out.Body.Environment = s.Environment
		if s.Matomo.URL != "" {
			out.Body.Matomo = &s.Matomo
		}
		return out, nil
	})
}
