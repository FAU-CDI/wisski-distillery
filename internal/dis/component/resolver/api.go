//spellchecker:words resolver
package resolver

//spellchecker:words context http github wisski distillery internal component auth scopes pkglib httpx
import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/api"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"go.tkw01536.de/pkglib/httpx"
)

type API struct {
	component.Base
	dependencies struct {
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
		Config: component.GetStill(a).Config,
		Auth:   a.dependencies.Auth,

		Methods: []string{"GET"},

		Scope: scopes.ScopeResolver,
		Handler: func(s string, r *http.Request) (string, error) {
			uri := r.URL.Query().Get("uri")
			if uri == "" {
				return "", httpx.ErrBadRequest
			}
			target := a.dependencies.Resolver.Target(uri)
			if target == "" {
				return "", httpx.ErrNotFound
			}
			return target, nil
		},
	}, nil
}
