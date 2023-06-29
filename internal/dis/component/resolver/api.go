package resolver

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/api"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/tkw1536/pkglib/httpx"
)

type API struct {
	component.Base
	Dependencies struct {
		Auth     *auth.Auth
		Resolver *Resolver
	}
}

var (
	_ component.Routeable = (*API)(nil)
)

func (api *API) Routes() component.Routes {
	return component.Routes{
		Prefix: "/api/v1/resolve/",
		Exact:  true,
	}
}

func (a *API) HandleRoute(ctx context.Context, path string) (http.Handler, error) {
	return &api.Handler[string]{
		Config: a.Config,
		Auth:   a.Dependencies.Auth,

		Methods: []string{"GET"},

		Scope: scopes.ScopeResolver,
		Handler: func(s string, r *http.Request) (string, error) {
			uri := r.URL.Query().Get("uri")
			if uri == "" {
				return "", httpx.ErrBadRequest
			}
			target := a.Dependencies.Resolver.Target(uri)
			if target == "" {
				return "", httpx.ErrNotFound
			}
			return target, nil
		},
	}, nil
}
