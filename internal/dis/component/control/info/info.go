package info

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter/logger"
	"github.com/gorilla/mux"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
	"github.com/tkw1536/goprogram/stream"
)

type Info struct {
	component.Base

	Analytics *lazy.PoolAnalytics
	Fetchers  []component.DistilleryFetcher

	Exporter     *exporter.Exporter
	Instances    *instances.Instances
	SnapshotsLog *logger.Logger
}

func (*Info) Routes() []string { return []string{"/dis/"} }

func (info *Info) Handler(route string, context context.Context, io stream.IOStream) (handler http.Handler, err error) {
	router := mux.NewRouter()
	{
		socket := &httpx.WebSocket{
			Context:  context,
			Fallback: router,
			Handler:  info.serveSocket,
		}
		handler = httpx.BasicAuth(socket, "WissKI Distillery Admin", func(user, pass string) bool {
			return user == info.Config.DisAdminUser && pass == info.Config.DisAdminPassword
		})
	}

	// handle everything
	router.Path(route).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, route+"/index", http.StatusTemporaryRedirect)
	})

	// add a handler for the index page
	router.Path(route + "index").Handler(httpx.HTMLHandler[indexContext]{
		Handler:  info.index,
		Template: indexTemplate,
	})

	// add a handler for the component page
	router.Path(route + "components").Handler(httpx.HTMLHandler[componentContext]{
		Handler:  info.components,
		Template: componentsTemplate,
	})

	// add a handler for the component page
	router.Path(route + "ingredients/{slug}").Handler(httpx.HTMLHandler[ingredientsContext]{
		Handler:  info.ingredients,
		Template: ingredientsTemplate,
	})

	// add a handler for the instance page
	router.Path(route + "instance/{slug}").Handler(httpx.HTMLHandler[instanceContext]{
		Handler:  info.instance,
		Template: instanceTemplate,
	})

	router.Path(route + "api/login").Handler(httpx.ClientSideRedirect(func(r *http.Request) (string, error) {
		// enforce POST
		if r.Method != http.MethodPost {
			return "", httpx.ErrMethodNotAllowed
		}

		// parse the form
		if err := r.ParseForm(); err != nil {
			return "", err
		}

		// get the instance
		instance, err := info.Instances.WissKI(r.PostFormValue("slug"))
		if err != nil {
			return "", httpx.ErrNotFound
		}

		target, err := instance.Users().Login(nil, r.PostFormValue("user"))
		if err != nil {
			return "", err
		}
		return target.String(), err
	}))

	return
}
