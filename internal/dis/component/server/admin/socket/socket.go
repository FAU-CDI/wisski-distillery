package socket

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/process_over_websocket"
	"github.com/FAU-CDI/process_over_websocket/proto"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances/purger"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/provision"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/admin/socket/actions"
	"github.com/tkw1536/pkglib/lazy"
)

type Sockets struct {
	component.Base

	handler lazy.Lazy[proto.Handler]

	dependencies struct {
		Actions  []actions.WebsocketAction
		IActions []actions.WebsocketInstanceAction

		Provision *provision.Provision
		Instances *instances.Instances
		Exporter  *exporter.Exporter
		Purger    *purger.Purger
		Auth      *auth.Auth
	}
}

var (
	_ component.Routeable = (*Sockets)(nil)
)

func (socket *Sockets) Routes() component.Routes {
	return component.Routes{
		Prefix:    "/api/v1/pow",
		Exact:     true,
		Decorator: socket.dependencies.Auth.Require(true, scopes.ScopeUserValid, nil),
	}
}

func (sockets *Sockets) HandleRoute(ctx context.Context, path string) (http.Handler, error) {
	server := process_over_websocket.Server{
		Handler: sockets.handler.Get(func() proto.Handler { return sockets.Actions(ctx) }),
		Options: process_over_websocket.Options{
			DisableREST: true, // for now
		},
	}

	// ensure that the server is closed once we are
	go func() {
		<-ctx.Done()
		server.Close()
	}()

	return http.StripPrefix("/api/v1/pow", &server), nil
}
