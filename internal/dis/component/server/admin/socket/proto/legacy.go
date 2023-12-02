package proto

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/tkw1536/pkglib/httpx/websocket"
	"github.com/tkw1536/pkglib/recovery"
)

var (
	errLegacyReadParamsTimeout = errors.New("timeout reading the first message")
	errLegacyUnknownAction     = errors.New("unknown action call")
	errLegacyIncorrectParams   = errors.New("invalid number of parameters")
)

// Handle handles the legacy protocol version.
// This is mostly used for legacy clients.
//
// There are two kinds of messages:
//
// - text messages, which are used to send input and output.
// - binary messages, which are json-encoded and used for control flow.
//
// To call an action, a client should send a [LegacyCallMessage] struct.
// The server will then start handling input and output (via text messages).
// If the client sends a SignalMessage, the signal is propagnated to the underlying context.
// Finally it will send a ResultMessage once handling is complete.
//
// A corresponding client implementation of this can be found in ..../remote/proto.ts
func (am ActionMap) handleLegacyProtocol(auth *auth.Auth, conn *websocket.Connection) (name string, err error) {
	var wg sync.WaitGroup

	// once we have finished executing send a binary message (indicating success) to the client.
	defer func() {
		// close the underlying connection, and then wait for everything to finish!
		defer wg.Wait()
		defer conn.Close()

		// recover from any errors
		if e := recovery.Recover(recover()); e != nil {
			err = e
		}

		// generate a result message
		var result LegacyResultMessage
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
		var message websocket.Message
		message.Type = websocket.BinaryMessage
		message.Bytes, err = json.Marshal(result)

		// silently fail if the message fails to encode
		// although this should not happen
		if err != nil {
			return
		}

		// and tell the client about it!
		<-conn.Write(message)
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
					textMessages <- string(msg.Bytes)
				}
				if msg.Type == websocket.BinaryMessage {
					binaryMessages <- msg.Bytes
				}
			case <-conn.Context().Done():
				return
			}
		}

	}()

	var call LegacyCallMessage
	select {
	case buffer := <-binaryMessages:
		if err := json.Unmarshal(buffer, &call); err != nil {
			return "", errLegacyUnknownAction
		}

	case <-time.After(1 * time.Second):
		return "", errLegacyReadParamsTimeout
	}

	// check that the given action exists!
	// and has the right number of parameters!
	action, ok := am[call.Call]
	if !ok || action.Handle == nil {
		return call.Call, errLegacyUnknownAction
	}
	if action.NumParams != len(call.Params) {
		return call.Call, errLegacyIncorrectParams
	}

	// check that we have the given permission
	if err := auth.CheckScope(action.ScopeParam, action.scope(), conn.Request()); err != nil {
		return call.Call, err
	}

	// create a context to be canceled once done
	ctx, cancel := context.WithCancel(conn.Context())
	defer cancel()

	// handle any signal messages
	wg.Add(1)
	go func() {
		defer wg.Done()
		var signal LegacySignalMessage

		for binary := range binaryMessages {
			signal.Signal = ""

			// read the signal message
			if err := json.Unmarshal(binary, &signal); err != nil {
				continue
			}

			// if we got a cancel message, do the cancellation!
			if signal.Signal == LegacySignalCancel {
				cancel()
			}
		}
	}()

	// create a pipe to handle the input
	// and start handling it
	var inputR, inputW = io.Pipe()
	defer inputW.Close()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for text := range textMessages {
			inputW.Write([]byte(text))
		}
	}()

	// write the output to the client as it comes in!
	// NOTE(twiesing): We may eventually need buffering here ...
	output := WriterFunc(func(b []byte) (int, error) {
		<-conn.WriteText(string(b))
		return len(b), nil
	})

	// handle the actual
	return call.Call, action.Handle(ctx, inputR, output, call.Params...)
}

// LegacyCallMessage is sent by the client to the server to invoke a remote procedure
type LegacyCallMessage struct {
	Call   string   `json:"call"`
	Params []string `json:"params,omitempty"`
}

// LegacySignalMessage is sent from the client to the server to stop the current procedure
type LegacySignalMessage struct {
	Signal LegacySignal `json:"signal"`
}

type LegacySignal string

const (
	LegacySignalCancel LegacySignal = "cancel"
)

// LegacyResultMessage is sent by the server to the client to report the success of a remote procedure
type LegacyResultMessage struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
