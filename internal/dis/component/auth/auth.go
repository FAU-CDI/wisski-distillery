package auth

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

type Auth struct {
	component.Base
	Dependencies struct {
		SQL             *sql.SQL
		UserDeleteHooks []component.UserDeleteHook
	}

	store lazy.Lazy[sessions.Store]
	csrf  lazy.Lazy[func(http.Handler) http.Handler]
}

var (
	_ component.Routeable = (*Auth)(nil)
)

func (auth *Auth) Routes() []string { return []string{"/user/"} }

func (auth *Auth) HandleRoute(ctx context.Context, route string) (http.Handler, error) {
	router := httprouter.New()

	router.Handler(http.MethodGet, route, auth.authUser(ctx))

	{
		login := auth.authLogin(ctx)
		router.Handler(http.MethodGet, route+"login", login)
		router.Handler(http.MethodPost, route+"login", login)
	}

	router.Handler(http.MethodGet, route+"logout", auth.authLogout(ctx))

	{
		password := auth.authPassword(ctx)
		router.Handler(http.MethodGet, route+"password", password)
		router.Handler(http.MethodPost, route+"password", password)
	}

	{
		totpenable := auth.authTOTPEnable(ctx)
		router.Handler(http.MethodGet, route+"totp/enable", totpenable)
		router.Handler(http.MethodPost, route+"totp/enable", totpenable)
	}

	{
		totpenroll := auth.authTOTPEnroll(ctx)
		router.Handler(http.MethodGet, route+"totp/enroll", totpenroll)
		router.Handler(http.MethodPost, route+"totp/enroll", totpenroll)
	}

	{
		totpdisable := auth.authTOTPDisable(ctx)
		router.Handler(http.MethodGet, route+"totp/disable", totpdisable)
		router.Handler(http.MethodPost, route+"totp/disable", totpdisable)
	}

	return router, nil
}

func (auth *Auth) CSRF() func(http.Handler) http.Handler {
	// setup the csrf handler (if needed)
	return auth.csrf.Get(func() func(http.Handler) http.Handler {
		var opts []csrf.Option
		if !auth.Config.HTTPSEnabled() {
			opts = append(opts, csrf.Secure(false))
		}
		return csrf.Protect(auth.Config.CSRFSecret(), opts...)
	})
}
