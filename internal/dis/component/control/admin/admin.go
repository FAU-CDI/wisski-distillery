package admin

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/policy"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static/custom"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter/logger"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances/purger"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
)

type Admin struct {
	component.Base
	Dependencies struct {
		Fetchers []component.DistilleryFetcher

		Exporter     *exporter.Exporter
		Instances    *instances.Instances
		SnapshotsLog *logger.Logger

		Auth *auth.Auth

		Policy *policy.Policy

		Custom *custom.Custom

		Purger *purger.Purger
	}

	Analytics *lazy.PoolAnalytics
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

func (admin *Admin) HandleRoute(ctx context.Context, route string) (handler http.Handler, err error) {

	router := httprouter.New()

	{
		handler = &httpx.WebSocket{
			Context:  ctx,
			Fallback: router,
			Handler:  admin.serveSocket,
		}
	}

	// add a handler for the index page
	router.Handler(http.MethodGet, route, httpx.HTMLHandler[indexContext]{
		Handler:  admin.index,
		Template: admin.Dependencies.Custom.Template(indexTemplate),
	})

	// fallback to the "/" page
	router.HandlerFunc(http.MethodGet, route+"index", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, route, http.StatusTemporaryRedirect)
	})

	// add a handler for the user page
	router.Handler(http.MethodGet, route+"users", httpx.HTMLHandler[userContext]{
		Handler:  admin.users,
		Template: admin.Dependencies.Custom.Template(userTemplate),
	})

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

	// add a handler for the component page
	router.Handler(http.MethodGet, route+"components", httpx.HTMLHandler[componentContext]{
		Handler:  admin.components,
		Template: admin.Dependencies.Custom.Template(componentsTemplate),
	})

	// add a handler for the component page
	router.Handler(http.MethodGet, route+"ingredients/:slug", httpx.HTMLHandler[ingredientsContext]{
		Handler:  admin.ingredients,
		Template: admin.Dependencies.Custom.Template(ingredientsTemplate),
	})

	// add a handler for the instance page
	router.Handler(http.MethodGet, route+"instance/:slug", httpx.HTMLHandler[instanceContext]{
		Handler:  admin.instance,
		Template: admin.Dependencies.Custom.Template(instanceTemplate),
	})

	// add a router for the grants pages
	router.Handler(http.MethodGet, route+"grants/:slug", httpx.HTMLHandler[grantsContext]{
		Handler:  admin.getGrants,
		Template: admin.Dependencies.Custom.Template(grantsTemplate),
	})
	router.Handler(http.MethodPost, route+"grants/", httpx.HTMLHandler[grantsContext]{
		Handler:  admin.postGrants,
		Template: admin.Dependencies.Custom.Template(grantsTemplate),
	})

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
