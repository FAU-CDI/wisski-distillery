//spellchecker:words proto
package proto

//spellchecker:words errors github wisski distillery internal component auth pkglib errorsx websocketx
import (
	"errors"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/websocketx"
)

// ActionMap handles a set of WebSocket actions.
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
		errWrite := conn.WritePrepared(msgUnknownSubprotocol)
		if errWrite != nil {
			errWrite = fmt.Errorf("unable to report unknown subprotocol to client: %w", errWrite)
		}
		return "", errorsx.Combine(errUnknownSubprotocol, errWrite)
	}
}

// WriterFunc implements io.Writer using a function.
type WriterFunc func([]byte) (int, error)

func (wf WriterFunc) Write(b []byte) (int, error) {
	return wf(b)
}
