//spellchecker:words news
package news

//spellchecker:words context http github wisski distillery internal component auth scopes
import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/api"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
)

type API struct {
	component.Base
	dependencies struct {
		Auth *auth.Auth
	}
}

var (
	_ component.Routeable = (*API)(nil)
)

func (api *API) Routes() component.Routes {
	return component.Routes{
		Prefix: "/api/v1/news/",
		Exact:  true,
	}
}

func (a *API) HandleRoute(ctx context.Context, path string) (http.Handler, error) {
	return &api.Handler[[]Item]{
		Config: component.GetStill(a).Config,
		Auth:   a.dependencies.Auth,

		Methods: []string{"GET"},

		Scope: scopes.ScopeListNews,
		Handler: func(s string, r *http.Request) ([]Item, error) {
			return Items()
		},
	}, nil
}
