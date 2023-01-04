package admin

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter/logger"
	"github.com/julienschmidt/httprouter"

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
	}

	Analytics *lazy.PoolAnalytics
}

var (
	_ component.DistilleryFetcher = (*Admin)(nil)
	_ component.Routeable         = (*Admin)(nil)
)

func (*Admin) Routes() []string { return []string{"/admin/"} }

func (admin *Admin) HandleRoute(ctx context.Context, route string) (handler http.Handler, err error) {

	router := httprouter.New()

	{
		socket := &httpx.WebSocket{
			Context:  ctx,
			Fallback: router,
			Handler:  admin.serveSocket,
		}
		handler = admin.Dependencies.Auth.Protect(socket, auth.Admin)
	}

	// handle everything
	router.HandlerFunc(http.MethodGet, route, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, route+"index", http.StatusTemporaryRedirect)
	})

	// add a handler for the index page
	router.Handler(http.MethodGet, route+"index", httpx.HTMLHandler[indexContext]{
		Handler:  admin.index,
		Template: indexTemplate,
	})

	// add a handler for the component page
	router.Handler(http.MethodGet, route+"components", httpx.HTMLHandler[componentContext]{
		Handler:  admin.components,
		Template: componentsTemplate,
	})

	// add a handler for the component page
	router.Handler(http.MethodGet, route+"ingredients/:slug", httpx.HTMLHandler[ingredientsContext]{
		Handler:  admin.ingredients,
		Template: ingredientsTemplate,
	})

	// add a handler for the instance page
	router.Handler(http.MethodGet, route+"instance/:slug", httpx.HTMLHandler[instanceContext]{
		Handler:  admin.instance,
		Template: instanceTemplate,
	})

	router.Handler(http.MethodPost, route+"api/login", httpx.RedirectHandler(func(r *http.Request) (string, int, error) {
		// parse the form
		if err := r.ParseForm(); err != nil {
			return "", 0, err
		}

		// get the instance
		instance, err := admin.Dependencies.Instances.WissKI(r.Context(), r.PostFormValue("slug"))
		if err != nil {
			return "", 0, httpx.ErrNotFound
		}

		target, err := instance.Users().Login(r.Context(), nil, r.PostFormValue("user"))
		if err != nil {
			return "", 0, err
		}
		return target.String(), http.StatusSeeOther, err
	}))

	return
}
