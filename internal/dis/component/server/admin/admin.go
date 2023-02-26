package admin

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/policy"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/admin/socket"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"github.com/tkw1536/pkglib/lifetime"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
)

type Admin struct {
	component.Base
	Dependencies struct {
		Fetchers []component.DistilleryFetcher

		Instances *instances.Instances

		Auth *auth.Auth

		Policy *policy.Policy

		Templating *templating.Templating

		Sockets *socket.Sockets
	}

	Analytics *lifetime.Analytics
}

var (
	_ component.DistilleryFetcher = (*Admin)(nil)
	_ component.Routeable         = (*Admin)(nil)
	_ component.Menuable          = (*Admin)(nil)
)

func (admin *Admin) Routes() component.Routes {
	return component.Routes{
		Prefix:    "/admin/",
		CSRF:      true,
		Decorator: admin.Dependencies.Auth.Require(auth.Admin),
	}
}

func (admin *Admin) Menu(r *http.Request) []component.MenuItem {
	if !admin.Dependencies.Auth.Has(auth.Admin, r) {
		return nil
	}
	return []component.MenuItem{
		{
			Title:    "Admin",
			Path:     "/admin/",
			Priority: component.MenuAdmin,
		},
	}
}

var (
	menuAdmin      = component.MenuItem{Title: "Admin", Path: "/admin/"}
	menuComponents = component.MenuItem{Title: "Components", Path: "/admin/components/", Priority: component.SmallButton}

	menuUsers      = component.MenuItem{Title: "Users", Path: "/admin/users/"}
	menuUserCreate = component.MenuItem{Title: "Create User", Path: "/admin/users/create/"}

	menuInstances   = component.MenuItem{Title: "Instances", Path: "/admin/instance/"}
	menuInstance    = component.DummyMenuItem()
	menuGrants      = component.DummyMenuItem()
	menuIngredients = component.DummyMenuItem()
)

func (admin *Admin) HandleRoute(ctx context.Context, route string) (handler http.Handler, err error) {

	router := httprouter.New()

	{
		handler = &httpx.WebSocket{
			Context:  ctx,
			Fallback: router,
			Handler:  admin.Dependencies.Sockets.Serve,
		}
	}

	// add a handler for the index page
	{
		index := admin.index(ctx)
		router.Handler(http.MethodGet, route, index)
	}

	// add a handler for the instances page
	{
		instances := admin.instances(ctx)
		router.Handler(http.MethodGet, route+"instance/", instances)
	}

	// add a handler for the user page
	{
		users := admin.users(ctx)
		router.Handler(http.MethodGet, route+"users", users)
	}

	// add a user create form
	{
		create := admin.createUser(ctx)
		router.Handler(http.MethodGet, route+"users/create", create)
		router.Handler(http.MethodPost, route+"users/create", create)
	}

	// add all the admin actions
	router.Handler(http.MethodPost, route+"users/delete", admin.usersDeleteHandler(ctx))
	router.Handler(http.MethodPost, route+"users/disable", admin.usersDisableHandler(ctx))
	router.Handler(http.MethodPost, route+"users/disabletotp", admin.usersDisableTOTPHandler(ctx))
	router.Handler(http.MethodPost, route+"users/password", admin.usersPasswordHandler(ctx))
	router.Handler(http.MethodPost, route+"users/toggleadmin", admin.usersToggleAdmin(ctx))
	router.Handler(http.MethodPost, route+"users/impersonate", admin.usersImpersonateHandler(ctx))
	router.Handler(http.MethodPost, route+"users/unsetpassword", admin.usersUnsetPasswordHandler(ctx))

	// add a handler for the component page
	{
		components := admin.components(ctx)
		router.Handler(http.MethodGet, route+"components", components)
	}

	// add a handler for the ingredients page
	{
		ingredients := admin.ingredients(ctx)
		router.Handler(http.MethodGet, route+"ingredients/:slug", ingredients)
	}

	// add a handler for the instance page
	{
		instance := admin.instance(ctx)
		router.Handler(http.MethodGet, route+"instance/:slug", instance)
	}

	{
		grants := admin.grants(ctx)
		router.Handler(http.MethodGet, route+"grants/:slug", grants)
		router.Handler(http.MethodPost, route+"grants/", grants) // NOTE(twiesing): This path is intentionally different!
	}

	// add a router for the login page
	router.Handler(http.MethodPost, route+"login", admin.loginHandler(ctx))

	return
}

func (admin *Admin) loginHandler(ctx context.Context) http.Handler {
	logger := zerolog.Ctx(ctx)

	return httpx.RedirectHandler(func(r *http.Request) (string, int, error) {
		// parse the form
		if err := r.ParseForm(); err != nil {
			logger.Err(err).Msg("failed to parse admin login")
			return "", 0, err
		}

		// get the instance
		instance, err := admin.Dependencies.Instances.WissKI(r.Context(), r.PostFormValue("slug"))
		if err != nil {
			return "", 0, httpx.ErrNotFound
		}

		target, err := instance.Users().Login(r.Context(), nil, r.PostFormValue("user"))
		if err != nil {
			logger.Err(err).Msg("failed to admin login")
			return "", 0, err
		}
		return target.String(), http.StatusSeeOther, err
	})
}
