//spellchecker:words proto
package proto

//spellchecker:words context encoding json errors sync time github wisski distillery internal component auth gorilla websocket pkglib errorsx recovery websocketx
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/gorilla/websocket"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/recovery"
	"go.tkw01536.de/pkglib/websocketx"
)

var (
	errReadParamsTimeout = errors.New("timeout reading the first message")
	errUnknownAction     = errors.New("unknown action call")
	errIncorrectParams   = errors.New("invalid number of parameters")
)

// A corresponding client implementation of this can be found in ..../remote/proto.ts.
func (am ActionMap) handleV1Protocol(auth *auth.Auth, conn *websocketx.Connection) (name string, err error) {
	var wg sync.WaitGroup

	// once we have finished executing send a binary message (indicating success) to the client.
	defer func() {
		// close the underlying connection, and then wait for everything to finish!
		defer wg.Wait()
		defer errorsx.Close(conn, &err, "connection")

		// recover from any errors
		if e := recovery.Recover(recover()); e != nil {
			err = e
		}

		// generate a result message
		var result ResultMessage
		if err == nil {
			result.Success = true
		} else {
			result.Success = false
			result.Message = err.Error()
			if result.Message == "" {
				result.Message = "unspecified error"
			}
		}

		// encode the result message to json!
		message := websocketx.NewBinaryMessage(nil)
		message.Body, err = json.Marshal(result)

		// silently fail if the message fails to encode
		// although this should not happen
		if err != nil {
			return
		}

		// and tell the client about it!
		if e := conn.Write(message); e != nil {
			e = fmt.Errorf("failed to write result message: %w", e)
			err = errorsx.Combine(err, e)
		}
	}()

	// create channels to receive text and bytes messages
	textMessages := make(chan string, 10)
	binaryMessages := make(chan []byte, 10)

	// start reading text and binary messages
	// and redirect everything to the right channels
	wg.Add(1)
	go func() {
		defer wg.Done()

		defer close(textMessages)
		defer close(binaryMessages)

		for {
			select {
			case msg := <-conn.Read():
				if msg.Type == websocket.TextMessage {
					textMessages <- string(msg.Body)
				}
				if msg.Type == websocket.BinaryMessage {
					binaryMessages <- msg.Body
				}
			case <-conn.Context().Done():
				return
			}
		}
	}()

	var call CallMessage
	select {
	case buffer := <-binaryMessages:
		if err := json.Unmarshal(buffer, &call); err != nil {
			return "", errUnknownAction
		}

	case <-time.After(1 * time.Second):
		return "", errReadParamsTimeout
	}

	// check that the given action exists!
	// and has the right number of parameters!
	action, ok := am[call.Call]
	if !ok || action.Handle == nil {
		return call.Call, errUnknownAction
	}
	if action.NumParams != len(call.Params) {
		return call.Call, errIncorrectParams
	}

	// check that we have the given permission
	if err := auth.CheckScope(action.ScopeParam, action.scope(), conn.Request()); err != nil {
		return call.Call, fmt.Errorf("failed to check scope: %w", err)
	}

	// create a context to be canceled once done
	ctx, cancel := context.WithCancel(conn.Context())
	defer cancel()

	// handle any signal messages
	wg.Add(1)
	go func() {
		defer wg.Done()
		var signal SignalMessage

		for binary := range binaryMessages {
			signal.Signal = ""

			// read the signal message
			if err := json.Unmarshal(binary, &signal); err != nil {
				continue
			}

			// if we got a cancel message, do the cancellation!
			if signal.Signal == SignalCancel {
				cancel()
			}
		}
	}()

	// create a pipe to handle the input
	// and start handling it
	var inputR, inputW = io.Pipe()
	defer errorsx.Close(inputW, &err, "input writer")

	wg.Add(1)
	go func() {
		defer wg.Done()

		for text := range textMessages {
			_, _ = inputW.Write([]byte(text)) // no way to report this error
		}
	}()

	// write the output to the client as it comes in!
	// NOTE(twiesing): We may eventually need buffering here ...
	output := WriterFunc(func(b []byte) (int, error) {
		if err := conn.WriteText(string(b)); err != nil {
			return 0, fmt.Errorf("failed to write text: %w", err)
		}
		return len(b), nil
	})

	// handle the actual
	return call.Call, action.Handle(ctx, inputR, output, call.Params...)
}

// CallMessage is sent by the client to the server to invoke a remote procedure.
type CallMessage struct {
	Call   string   `json:"call"`
	Params []string `json:"params,omitempty"`
}

// SignalMessage is sent from the client to the server to stop the current procedure.
type SignalMessage struct {
	Signal Signal `json:"signal"`
}

type Signal string

const (
	SignalCancel Signal = "cancel"
)

// ResultMessage is sent by the server to the client to report the success of a remote procedure.
type ResultMessage struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
