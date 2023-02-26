package phpx

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"

	_ "embed"

	"github.com/tkw1536/goprogram/stream"
	"github.com/tkw1536/pkglib/collection"
	"github.com/tkw1536/pkglib/contextx"
	"github.com/tkw1536/pkglib/lazy"
	"github.com/tkw1536/pkglib/nobufio"
)

// Server represents a server that executes PHP code.
// A typical use-case is to define functions using [MarshalEval], and then call those functions [MarshalCall].
//
// A server, once used, should be closed using the [Close] method.
type Server struct {
	// Context to use for the server
	Context context.Context

	// Executor is the executor used by this server.
	// It may not be modified concurrently with other processes.
	Executor Executor

	// prepares the server
	init sync.Once
	err  lazy.Lazy[error]

	// input / output for underlying executor
	in  io.WriteCloser
	out io.Reader

	m sync.Mutex // prevents concurrent access on any of the methods

	cancel context.CancelFunc
	c      context.Context // closed when server is finished
}

func (server *Server) prepare() error {
	server.init.Do(func() {
		// create input and output pipes
		ir, iw, err := os.Pipe()
		if err != nil {
			server.err.Set(ServerError{errInit, err})
			return
		}
		or, ow, err := os.Pipe()
		if err != nil {
			ir.Close()
			iw.Close()

			server.err.Set(ServerError{errInit, err})
			return
		}

		// create a context to close the server
		context, cancel := context.WithCancel(server.Context)
		server.cancel = cancel

		// start the shell process, which will close everything once done
		go func() {
			defer func() {
				ir.Close()
				iw.Close()
				or.Close()
				ow.Close()

				server.cancel()
			}()

			// start the server
			io := stream.NewIOStream(ow, nil, ir, 0)
			err := server.Executor.Spawn(server.c, io, serverPHP)
			server.err.Set(ServerError{errClosed, err})
		}()

		server.in = iw
		server.out = or
		server.c = context
	})

	return server.err.Get(nil)
}

// MarshalEval evaluates code on the server and Marshals the result into value.
// When value is nil, the results are discarded.
//
// code is directly passed to php's "eval" function.
// as such any functions defined will remain in server memory.
//
// When an exception is thrown by the PHP Code, error is not nil, and dest remains unchanged.
func (server *Server) MarshalEval(ctx context.Context, value any, code string) error {
	if err := server.prepare(); err != nil {
		return err
	}

	// input to be sent to the server
	delim := findDelimiter(code)
	input := delim + "\n" + code + "\n" + delim + "\n"

	server.m.Lock()
	defer server.m.Unlock()

	// when the server is already done
	if err := server.c.Err(); err != nil {
		return ServerError{Message: errClosed}
	}

	// find a delimiter for the code, and then send
	io.WriteString(server.in, input)

	data, err, _ := contextx.Run2(ctx, func(start func()) (string, error) {
		return nobufio.ReadLine(server.out)
	}, func() {
		server.cancel()
	})
	if err != nil {
		return ServerError{Message: errReceive, Err: err}
	}

	// read whatever we received
	var received [2]json.RawMessage
	if err := json.Unmarshal([]byte(data), &received); err != nil {
		return ServerError{Message: errReceive, Err: err}
	}

	// check if there was an error
	var errString string
	if err := json.Unmarshal(received[1], &errString); err == nil && errString != "" {
		return Throwable(errString)
	}

	// special case: no return value => no unmarshaling needed
	if value == nil {
		return nil
	}

	// read the actual result!
	return json.Unmarshal(received[0], value)
}

// Eval is like [MarshalEval], but returns the value as an any
func (server *Server) Eval(ctx context.Context, code string) (value any, err error) {
	err = server.MarshalEval(ctx, &value, code)
	return
}

// MarshalCall calls a previously defined function with the given arguments.
// Arguments are sent to php using json Marshal, and are 'json_decode'd on the php side.
//
// Return values are received as in [MarshalEval].
func (server *Server) MarshalCall(ctx context.Context, value any, function string, args ...any) error {
	// name of function to call
	name := MarshalString(function)

	// generate code to call
	var code string
	switch len(args) {
	case 0:
		code = "return call_user_func(" + name + ");"
	case 1:
		param, err := Marshal(args[0])
		if err != nil {
			return err
		}
		code = "return call_user_func(" + name + "," + param + ");"
	default:
		params, err := Marshal(args)
		if err != nil {
			return err
		}
		code = "return call_user_func_array(" + name + "," + params + ");"
	}

	// and evaluate the code
	return server.MarshalEval(ctx, value, code)
}

// Call is like [MarshalCall] but returns the return value of the function as an any
func (server *Server) Call(ctx context.Context, function string, args ...any) (value any, err error) {
	err = server.MarshalCall(ctx, &value, function, args...)
	return
}

const delimiterRune = 'F' // press to pay respect

// findDelimiter finds a delimiter that does not occur in the input string
func findDelimiter(input string) string {
	// find the longest sequence of delimiter rune
	var current, longest int
	for _, r := range input {
		if r == delimiterRune {
			current++
		} else {
			current = 0
		}

		if current > longest {
			longest = current
		}
	}
	// and then return it multipled longer than that
	return strings.Repeat(string(delimiterRune), longest+1)
}

// Close closes this server and prevents any further code from being run.
func (server *Server) Close() error {
	server.prepare()

	server.m.Lock()
	defer server.m.Unlock()

	// if the context is already closed
	if err := server.c.Err(); err != nil {
		return ServerError{Message: errClosed}
	}

	server.in.Close()
	<-server.c.Done()

	return nil
}

//go:embed server.php
var serverPHP string

// pre-process the server.php code to make it shorter
func init() {
	minifier := regexp.MustCompile(`\s*([=)(.,{}])\s*`)

	// remove the first '<?php' line
	lines := strings.Split(serverPHP, "\n")[1:]
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	// remove comment lines
	lines = collection.Filter(lines, func(line string) bool {
		return !strings.HasPrefix(line, "//")
	})

	serverPHP = minifier.ReplaceAllString(strings.Join(lines, ""), "$1")
}
