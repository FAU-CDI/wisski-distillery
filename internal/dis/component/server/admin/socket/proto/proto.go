package proto

import (
	"errors"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/tkw1536/pkglib/websocketx"
)

// ActionMap handles a set of WebSocket actions
type ActionMap map[string]Action

var (
	errUnknownSubprotocol = errors.New("unknown subprotocol")
	msgUnknownSubprotocol = websocketx.NewTextMessage(errUnknownSubprotocol.Error()).MustPrepare()
)

// Handle handles a new incoming websocket connection by switching on the subprotocol.
// See appropriate protocol handlers for documentation.
func (am ActionMap) Handle(auth *auth.Auth, conn *websocketx.Connection) (name string, err error) {
	// select based on the negotiated subprotocol
	switch conn.Subprotocol() {
	case "":
		return am.handleV1Protocol(auth, conn)
	default:
		conn.WritePrepared(msgUnknownSubprotocol)
		return "", errUnknownSubprotocol
	}
}

// WriterFunc implements io.Writer using a function.
type WriterFunc func([]byte) (int, error)

func (wf WriterFunc) Write(b []byte) (int, error) {
	return wf(b)
}
