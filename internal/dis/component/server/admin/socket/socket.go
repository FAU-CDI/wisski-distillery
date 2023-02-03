package socket

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances/purger"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/tkw1536/goprogram/status"
)

type Sockets struct {
	component.Base

	Dependencies struct {
		Instances *instances.Instances
		Exporter  *exporter.Exporter
		Purger    *purger.Purger
	}
}

// Serve handles a connection to the websocket api
func (socket *Sockets) Serve(conn httpx.WebSocketConnection) {
	// read the next message to act on
	message, ok := <-conn.Read()
	if !ok {
		return
	}

	name := string(message.Bytes)

	// perform a generic action first
	if action, ok := actions[name]; ok {
		socket.Handle(conn, action)
		return
	}

	// then do the socket actions
	if action, ok := igActions[name]; ok {
		socket.Handle(conn, action)
	}
}

var instanceParamsTimeout = time.Second

func (socket *Sockets) Handle(conn httpx.WebSocketConnection, action SocketAction) {
	// read the parameters
	params := make([]string, action.NumParams)
	for i := range params {
		select {
		case message, ok := <-conn.Read():
			if !ok {
				<-conn.WriteText("Insufficient parameters")
				return
			}
			params[i] = string(message.Bytes)
		case <-time.After(instanceParamsTimeout):
			<-conn.WriteText("Timed out reading parameters")
			return
		}
	}

	// build a stream
	writer := &status.LineBuffer{
		Line: func(line string) {
			<-conn.WriteText(line)
		},
		FlushLineOnClose: true,
	}
	defer writer.Close()

	// handle the interactive action
	if action.HandleInteractive != nil {
		err := action.HandleInteractive(conn.Context(), socket, writer, params...)
		if err != nil {
			fmt.Fprintln(writer, err)
			return
		}
		fmt.Fprintln(writer, "done")
	}
}

// IAction is like SocketAction, but takes the slug of an instance (runnning or not) as the first parameter
type IAction struct {
	NumParams int

	HandleInteractive func(ctx context.Context, sockets *Sockets, instance *wisski.WissKI, out io.Writer, params ...string) error
}

// AsGenericAction turns this InstanceAction into a generic action
func (ia IAction) AsGenericAction() SocketAction {
	return SocketAction{
		NumParams: ia.NumParams + 1,
		HandleInteractive: func(ctx context.Context, sockets *Sockets, out io.Writer, params ...string) error {
			instance, err := sockets.Dependencies.Instances.WissKI(ctx, params[0])
			if err != nil {
				return err
			}

			return ia.HandleInteractive(ctx, sockets, instance, out, params[1:]...)
		},
	}
}

// SocketAction represents an action handled via socket
type SocketAction struct {
	NumParams int

	HandleInteractive func(ctx context.Context, sockets *Sockets, out io.Writer, params ...string) error
}
