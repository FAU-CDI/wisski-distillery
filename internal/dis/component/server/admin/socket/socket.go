package socket

import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
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

	Dependencies struct {
		Provision *provision.Provision
		Instances *instances.Instances
		Exporter  *exporter.Exporter
		Purger    *purger.Purger
	}
}

// Serve handles a connection to the websocket api
func (socket *Sockets) Serve(conn httpx.WebSocketConnection) {
	// handle the websocket connection!
	name, err := socket.actions.Get(socket.Actions).Handle(conn)
	if err != nil {
		zerolog.Ctx(conn.Context()).Err(err).Str("name", name).Msg("Error handling websocket")
	}
}

// Generic returns a new action that calls handler with the provided number of parameters
func (sockets *Sockets) Generic(numParams int, handler func(ctx context.Context, socket *Sockets, in io.Reader, out io.Writer, params ...string) error) Action {
	return Action{
		NumParams: numParams,
		Handle: func(ctx context.Context, in io.Reader, out io.Writer, params ...string) error {
			return handler(ctx, sockets, in, out, params...)
		},
	}
}

// Insstance returns a new action that calls handler with a specific WissKI instance
func (sockets *Sockets) Instance(numParams int, handler func(ctx context.Context, sockets *Sockets, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error) Action {
	return Action{
		NumParams: numParams + 1,
		Handle: func(ctx context.Context, in io.Reader, out io.Writer, params ...string) error {
			instance, err := sockets.Dependencies.Instances.WissKI(ctx, params[0])
			if err != nil {
				return err
			}
			return handler(ctx, sockets, instance, in, out, params[1:]...)
		},
	}
}
