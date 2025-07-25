//spellchecker:words auth
package auth

//spellchecker:words context http github wisski distillery internal component auth tokens server templating gorilla sessions julienschmidt httprouter pkglib lazy
import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/tokens"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"go.tkw01536.de/pkglib/lazy"
)

type Auth struct {
	component.Base
	dependencies struct {
		SQL             *sql.SQL
		UserDeleteHooks []component.UserDeleteHook
		Templating      *templating.Templating
		ScopeProviders  []component.ScopeProvider
		Tokens          *tokens.Tokens
	}

	store lazy.Lazy[sessions.Store]

	scopeMap lazy.Lazy[map[component.Scope]scopeMapEntry]
}

type scopeMapEntry struct {
	Provider component.ScopeProvider
	Info     component.ScopeInfo
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
