package panel

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/next"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/policy"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static/custom"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/julienschmidt/httprouter"
)

type UserPanel struct {
	component.Base
	Dependencies struct {
		Auth      *auth.Auth
		Custom    *custom.Custom
		Policy    *policy.Policy
		Instances *instances.Instances
		Next      *next.Next
	}
}

var (
	_ component.Routeable = (*UserPanel)(nil)
)

func (panel *UserPanel) Routes() component.Routes {
	return component.Routes{
		Paths:     []string{"/user/"},
		CSRF:      true,
		Decorator: panel.Dependencies.Auth.Require(nil),
	}
}

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
	custom.BaseContext
	httpx.FormContext

	User *models.User
}

func (panel *UserPanel) UserFormContext(ctx httpx.FormContext, r *http.Request) any {
	user, err := panel.Dependencies.Auth.UserOf(r)

	uctx := userFormContext{FormContext: ctx}
	panel.Dependencies.Custom.Update(&uctx, r)
	if err == nil {
		uctx.User = &user.User
	}
	return uctx
}
