package api

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
)

type API struct {
	component.Base
	Dependencies struct {
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
		Config: a.Config,
		Auth:   a.Dependencies.Auth,

		Methods: []string{"GET"},

		Handler: func(s string, r *http.Request) (ai AuthInfo, err error) {
			var user *auth.AuthUser
			user, ai.Token, err = a.Dependencies.Auth.UserOf(r)
			if user != nil {
				ai.User = user.User.User
			}
			return
		},
	}, nil
}
