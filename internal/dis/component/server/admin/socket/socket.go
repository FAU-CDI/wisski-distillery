package socket

import (
	"context"
	"io"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances/purger"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/provision"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/rs/zerolog"
	"github.com/tkw1536/pkglib/httpx"
	"github.com/tkw1536/pkglib/lazy"
)

type Sockets struct {
	component.Base

	actions lazy.Lazy[ActionMap]

	dependencies struct {
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
	name, err := socket.actions.Get(socket.Actions).Handle(socket.dependencies.Auth, conn)
	if err != nil {
		zerolog.Ctx(conn.Context()).Err(err).Str("name", name).Msg("Error handling websocket")
	}
}

// Generic returns a new action that calls handler with the provided number of parameters
func (sockets *Sockets) Generic(scope component.Scope, scopeParam string, numParams int, handler func(ctx context.Context, socket *Sockets, in io.Reader, out io.Writer, params ...string) error) Action {
	return Action{
		Scope:      scope,
		ScopeParam: scopeParam,
		NumParams:  numParams,
		Handle: func(ctx context.Context, in io.Reader, out io.Writer, params ...string) error {
			return handler(ctx, sockets, in, out, params...)
		},
	}
}

// Insstance returns a new action that calls handler with a specific WissKI instance
func (sockets *Sockets) Instance(scope component.Scope, scopeParam string, numParams int, handler func(ctx context.Context, sockets *Sockets, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error) Action {
	return Action{
		Scope:      scope,
		ScopeParam: scopeParam,

		NumParams: numParams + 1,
		Handle: func(ctx context.Context, in io.Reader, out io.Writer, params ...string) error {
			instance, err := sockets.dependencies.Instances.WissKI(ctx, params[0])
			if err != nil {
				return err
			}
			return handler(ctx, sockets, instance, in, out, params[1:]...)
		},
	}
}
