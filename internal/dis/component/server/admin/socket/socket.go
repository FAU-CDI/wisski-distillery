package socket

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances/purger"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/provision"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/gorilla/websocket"
	"github.com/tkw1536/goprogram/status"
	"github.com/tkw1536/pkglib/httpx"
)

type Sockets struct {
	component.Base

	Dependencies struct {
		Provision *provision.Provision
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

type actionResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func (*Sockets) reportErrorToClient(conn httpx.WebSocketConnection, err error) {
	// create an action result
	var result actionResult
	if err == nil {
		result.Success = true
	} else {
		result.Success = false
		result.Error = err.Error()
	}

	// marshal the result, ignoring any error silently
	data, err := json.Marshal(result)
	if err != nil {
		return
	}

	// and send it as a binary message to the client
	<-conn.Write(httpx.WebSocketMessage{Type: websocket.BinaryMessage, Bytes: data})
}

var errInsufficientParams = errors.New("insufficient parameters")
var errParameterTimeout = errors.New("timed out reading parameters")

func (socket *Sockets) Handle(conn httpx.WebSocketConnection, action SocketAction) (err error) {
	// report the error to the client
	defer func() {
		// NOTE: the closure is needed here!
		socket.reportErrorToClient(conn, err)
	}()

	// read the parameters
	params := make([]string, action.NumParams)
	for i := range params {
		select {
		case message, ok := <-conn.Read():
			if !ok {
				return errInsufficientParams
			}
			params[i] = string(message.Bytes)
		case <-time.After(instanceParamsTimeout):
			return errParameterTimeout
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
	return action.HandleInteractive(conn.Context(), socket, writer, params...)
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
