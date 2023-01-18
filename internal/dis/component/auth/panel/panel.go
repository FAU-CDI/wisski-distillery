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
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/ssh2/sshkeys"
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
		Keys      *sshkeys.SSHKeys
	}
}

var (
	_ component.Routeable = (*UserPanel)(nil)
	_ component.Menuable  = (*UserPanel)(nil)
)

func (panel *UserPanel) Routes() component.Routes {
	return component.Routes{
		Prefix:    "/user/",
		CSRF:      true,
		Decorator: panel.Dependencies.Auth.Require(nil),
	}
}

func (panel *UserPanel) Menu(r *http.Request) []component.MenuItem {
	title := "Login"

	user, err := panel.Dependencies.Auth.UserOf(r)
	if user != nil && err == nil {
		title = user.User.User
	}
	return []component.MenuItem{
		{Title: title, Priority: component.MenuUser, Path: "/user/"},
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

	{
		ssh := panel.sshRoute(ctx)
		router.Handler(http.MethodGet, route+"ssh", ssh)
	}

	{
		add := panel.sshAddRoute(ctx)
		router.Handler(http.MethodGet, route+"ssh/add", add)
		router.Handler(http.MethodPost, route+"ssh/add", add)
	}

	{
		delete := panel.sshDeleteRoute(ctx)
		router.Handler(http.MethodPost, route+"ssh/delete", delete)
	}

	// ensure that the user is logged in!
	return panel.Dependencies.Auth.Protect(router, nil), nil
}

type userFormContext struct {
	custom.BaseContext
	httpx.FormContext

	User *models.User
}

func (panel *UserPanel) UserFormContext2(tpl *custom.Template[userFormContext], last component.MenuItem, gaps ...custom.BaseContextGaps) func(ctx httpx.FormContext, r *http.Request) any {
	var g custom.BaseContextGaps
	if len(gaps) > 1 {
		panic("UserFormContext2: gaps must be of length 0 or 1")
	}
	if len(gaps) == 1 {
		g = gaps[0]
	}
	g.Crumbs = []component.MenuItem{
		{Title: "User", Path: "/user/"},
		last,
	}

	return custom.MappedHandler(tpl, func(ctx httpx.FormContext, r *http.Request) (userFormContext, custom.BaseContextGaps) {
		uctx := userFormContext{FormContext: ctx}
		if user, err := panel.Dependencies.Auth.UserOf(r); err == nil {
			uctx.User = &user.User
		}
		return uctx, g
	})
}
