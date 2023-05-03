package auth

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/tkw1536/pkglib/lazy"
)

type Auth struct {
	component.Base
	Dependencies struct {
		SQL             *sql.SQL
		UserDeleteHooks []component.UserDeleteHook
		Templating      *templating.Templating
		ScopeProviders  []component.ScopeProvider
	}

	store lazy.Lazy[sessions.Store]

	scopeInfos lazy.Lazy[[]component.ScopeInfo]
	scopeIndex lazy.Lazy[map[component.Scope]int]
}

var (
	_ component.Routeable = (*Auth)(nil)
	_ component.Menuable  = (*Auth)(nil)
	_ component.Table     = (*Auth)(nil)
)

func (auth *Auth) Routes() component.Routes {
	return component.Routes{
		Prefix: "/auth/",
		CSRF:   true,
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
