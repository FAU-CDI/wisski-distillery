package socket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tkw1536/goprogram/status"
	"github.com/tkw1536/pkglib/httpx"
)

// ActionMap handles a set of WebSocket actions
type ActionMap map[string]Action

var errReadParamsTimeout = errors.New("timeout reading the first message")
var errUnknownAction = errors.New("unknown action call")
var errIncorrectParams = errors.New("invalid number of parameters")

type errPanic struct{ value any }

func (err errPanic) Error() string {
	return fmt.Sprintf("fatal error: %v", err.value)
}

// Handle handles a new incoming websocket connection.
//
// There are two kinds of messages:
//
// - text messages, which are used to send input and output.
// - binary messages, which are json-encoded and used for control flow.
//
// To call an action, a client should send a CallMessage struct.
// The server will then start handling input and output (via text messages).
// If the client sends a SignalMessage, the signal is propagnated to the underlying context.
// Finally it will send a ResultMessage once handling is complete.
//
// A corresponding client implementation of this can be found in ..../remote/proto.ts
func (am ActionMap) Handle(conn httpx.WebSocketConnection) (name string, err error) {
	var wg sync.WaitGroup

	// once we have finished executing send a binary message (indicating success) to the client.
	defer func() {
		// close the underlying connection, and then wait for everything to finish!
		defer wg.Wait()
		defer conn.Close()

		// recover from any errors
		if v := recover(); v != nil {
			err = errPanic{value: v}
		}

		// generate a result message
		var result ResultMessage
		if err == nil {
			result.Success = true
		} else {
			result.Success = false
			result.Message = err.Error()
		}

		// encode the result message to json!
		var message httpx.WebSocketMessage
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
	action, ok := am[call.Name]
	if !ok || action.Handle == nil {
		return call.Name, errUnknownAction
	}
	if action.NumParams != len(call.Params) {
		return call.Name, errIncorrectParams
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
	defer inputW.Close()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for text := range textMessages {
			inputW.Write([]byte(text))
		}
	}()

	// create a linebuffer to write the output line by line
	output := &status.LineBuffer{
		Line: func(line string) {
			<-conn.WriteText(line)
		},
		FlushLineOnClose: true,
	}
	defer output.Close()

	// handle the actual
	return call.Name, action.Handle(ctx, inputR, output, call.Params...)
}

// Action is something that can be handled via a WebSocket connection.
type Action struct {
	// NumPara
	NumParams int

	// Handle handles this action.
	//
	// ctx is closed once the underlying connection is closed.
	// out is an io.Writer that is automatically sent to the client.
	// params holds exactly NumParams parameters.
	Handle func(ctx context.Context, in io.Reader, out io.Writer, params ...string) error
}

// CallMessage is sent by the client to the server to invoke a remote procedure
type CallMessage struct {
	Name   string   `json:"name"`
	Params []string `json:"params,omitempty"`
}

// CancelMessage is sent from the client to the server to stop the current procedure
type SignalMessage struct {
	Signal Signal `json:"signal"`
}

type Signal string

const (
	SignalCancel = "cancel"
)

// ResultMessage is sent by the server to the client to report the success of a remote procedure
type ResultMessage struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
