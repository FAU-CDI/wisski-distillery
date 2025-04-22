package api

//spellchecker:words context http github wisski distillery internal component auth
import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
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
		Prefix: "/api/v1/auth/",
		Exact:  true,
	}
}

type AuthInfo struct {
	// User returns the authenticated user.
	// If there is no user, contains the empty string.
	User string

	// Token indicates if the user is authenticated with a token.
	Token bool
}

func (a *API) HandleRoute(ctx context.Context, path string) (http.Handler, error) {
	return &Handler[AuthInfo]{
		Config: component.GetStill(a).Config,
		Auth:   a.dependencies.Auth,

		Methods: []string{"GET"},

		Handler: func(s string, r *http.Request) (ai AuthInfo, err error) {
			session, _, err := a.dependencies.Auth.SessionOf(r)
			ai.User = session.Username()
			ai.Token = session.Token
			return
		},
	}, nil
}
