package info

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/component/exporter/logger"

	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/tkw1536/goprogram/stream"
)

type Info struct {
	component.ComponentBase

	Exporter     *exporter.Exporter
	Instances    *instances.Instances
	SnapshotsLog *logger.Logger
}

func (Info) Name() string { return "control-info" }

func (*Info) Routes() []string { return []string{"/dis/"} }

func (info *Info) Handler(route string, context context.Context, io stream.IOStream) (http.Handler, error) {
	mux := http.NewServeMux()

	// handle everything
	mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == route {
			http.Redirect(w, r, route+"/index", http.StatusTemporaryRedirect)
			return
		}
		http.NotFound(w, r)
	})

	// add a handler for the index page
	mux.Handle(route+"index", httpx.HTMLHandler[indexPageContext]{
		Handler:  info.indexPageAPI,
		Template: indexTemplate,
	})

	// add a handler for the instance page
	mux.Handle(route+"instance/", httpx.HTMLHandler[instancePageContext]{
		Handler:  info.instancePageAPI,
		Template: instanceTemplate,
	})

	handler := &httpx.WebSocket{
		Context:  context,
		Fallback: mux,
		Handler:  info.serveSocket,
	}

	// ensure that everyone is logged in!
	return httpx.BasicAuth(handler, "WissKI Distillery Admin", func(user, pass string) bool {
		return user == info.Config.DisAdminUser && pass == info.Config.DisAdminPassword
	}), nil
}
