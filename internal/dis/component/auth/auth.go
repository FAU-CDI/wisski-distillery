package auth

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static/custom"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

type Auth struct {
	component.Base
	Dependencies struct {
		SQL             *sql.SQL
		UserDeleteHooks []component.UserDeleteHook
		Custom          *custom.Custom
	}

	store lazy.Lazy[sessions.Store]
}

var (
	_ component.Routeable = (*Auth)(nil)
)

func (auth *Auth) Routes() component.Routes {
	return component.Routes{
		Paths: []string{"/auth/"},
		CSRF:  true,
	}
}

func (auth *Auth) HandleRoute(ctx context.Context, route string) (http.Handler, error) {
	router := httprouter.New()
	{
		login := auth.authLogin(ctx)
		router.Handler(http.MethodGet, route+"login", login)
		router.Handler(http.MethodPost, route+"login", login)
	}

	router.Handler(http.MethodGet, route+"logout", auth.authLogout(ctx))

	return router, nil
}
