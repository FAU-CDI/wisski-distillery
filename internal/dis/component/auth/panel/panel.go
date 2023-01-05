package panel

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/julienschmidt/httprouter"
)

type UserPanel struct {
	component.Base
	Dependencies struct {
		Auth *auth.Auth
	}
}

var (
	_ component.Routeable = (*UserPanel)(nil)
)

func (panel *UserPanel) Routes() []string { return []string{"/user/"} }

func (panel *UserPanel) HandleRoute(ctx context.Context, route string) (http.Handler, error) {
	router := httprouter.New()

	{
		user := panel.routeUser(ctx)
		router.Handler(http.MethodGet, route, user)
	}

	{
		password := panel.routePassword(ctx)
		router.Handler(http.MethodGet, route+"password", password)
		router.Handler(http.MethodPost, route+"password", password)
	}

	{
		totpenable := panel.routeTOTPEnable(ctx)
		router.Handler(http.MethodGet, route+"totp/enable", totpenable)
		router.Handler(http.MethodPost, route+"totp/enable", totpenable)
	}

	{
		totpenroll := panel.routeTOTPEnroll(ctx)
		router.Handler(http.MethodGet, route+"totp/enroll", totpenroll)
		router.Handler(http.MethodPost, route+"totp/enroll", totpenroll)
	}

	{
		totpdisable := panel.routeTOTPDisable(ctx)
		router.Handler(http.MethodGet, route+"totp/disable", totpdisable)
		router.Handler(http.MethodPost, route+"totp/disable", totpdisable)
	}

	// ensure that the user is logged in!
	return panel.Dependencies.Auth.Protect(router, nil), nil
}

type userFormContext struct {
	httpx.FormContext
	User *models.User
}

func (panel *UserPanel) UserFormContext(ctx httpx.FormContext, r *http.Request) any {
	user, err := panel.Dependencies.Auth.UserOf(r)

	uctx := userFormContext{FormContext: ctx}
	if err == nil {
		uctx.User = &user.User
	}
	return uctx
}
