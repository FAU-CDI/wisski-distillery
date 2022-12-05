package info

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter/logger"
	"github.com/julienschmidt/httprouter"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
)

type Info struct {
	component.Base
	Dependencies struct {
		Fetchers []component.DistilleryFetcher

		Exporter     *exporter.Exporter
		Instances    *instances.Instances
		SnapshotsLog *logger.Logger
	}

	Analytics *lazy.PoolAnalytics
}

var (
	_ component.DistilleryFetcher = (*Info)(nil)
	_ component.Routeable         = (*Info)(nil)
)

func (*Info) Routes() []string { return []string{"/dis/"} }

func (info *Info) HandleRoute(ctx context.Context, route string) (handler http.Handler, err error) {

	router := httprouter.New()

	{
		socket := &httpx.WebSocket{
			Context:  ctx,
			Fallback: router,
			Handler:  info.serveSocket,
		}
		handler = httpx.BasicAuth(socket, "WissKI Distillery Admin", func(user, pass string) bool {
			return user == info.Config.DisAdminUser && pass == info.Config.DisAdminPassword
		})
	}

	// handle everything
	router.HandlerFunc(http.MethodGet, route, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, route+"index", http.StatusTemporaryRedirect)
	})

	// add a handler for the index page
	router.Handler(http.MethodGet, route+"index", httpx.HTMLHandler[indexContext]{
		Handler:  info.index,
		Template: indexTemplate,
	})

	// add a handler for the component page
	router.Handler(http.MethodGet, route+"components", httpx.HTMLHandler[componentContext]{
		Handler:  info.components,
		Template: componentsTemplate,
	})

	// add a handler for the component page
	router.Handler(http.MethodGet, route+"ingredients/:slug", httpx.HTMLHandler[ingredientsContext]{
		Handler:  info.ingredients,
		Template: ingredientsTemplate,
	})

	// add a handler for the instance page
	router.Handler(http.MethodGet, route+"instance/:slug", httpx.HTMLHandler[instanceContext]{
		Handler:  info.instance,
		Template: instanceTemplate,
	})

	router.Handler(http.MethodPost, route+"api/login", httpx.RedirectHandler(func(r *http.Request) (string, int, error) {
		// parse the form
		if err := r.ParseForm(); err != nil {
			return "", 0, err
		}

		// get the instance
		instance, err := info.Dependencies.Instances.WissKI(r.Context(), r.PostFormValue("slug"))
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
