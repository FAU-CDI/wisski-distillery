package socket

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances/purger"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/provision"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/admin/socket/actions"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/admin/socket/proto"
	"github.com/rs/zerolog"
	"github.com/tkw1536/pkglib/httpx"
	"github.com/tkw1536/pkglib/lazy"
)

type Sockets struct {
	component.Base

	actions lazy.Lazy[proto.ActionMap]

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
		Prefix:    "/api/v1/ws",
		Exact:     true,
		Decorator: socket.dependencies.Auth.Require(true, scopes.ScopeUserValid, nil),
	}
}

func (sockets *Sockets) HandleRoute(ctx context.Context, path string) (http.Handler, error) {
	return &httpx.WebSocket{
		Context: ctx,
		Handler: sockets.Serve,
	}, nil
}

// Serve handles a connection to the websocket api
func (socket *Sockets) Serve(conn httpx.WebSocketConnection) {
	// handle the websocket connection!
	name, err := socket.actions.Get(func() proto.ActionMap { return socket.Actions(conn.Context()) }).Handle(socket.dependencies.Auth, conn)
	if err != nil {
		zerolog.Ctx(conn.Context()).Err(err).Str("name", name).Msg("Error handling websocket")
	}
}
