package phpserver

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"
	"sync"

	_ "embed"

	"github.com/tkw1536/goprogram/lib/collection"
	"github.com/tkw1536/goprogram/lib/nobufio"
	"github.com/tkw1536/goprogram/stream"
)

// New creates a new server, with execPHP as a method to call a PHP Shell.
func New(execPHP func(str stream.IOStream, script string)) (*Server, error) {
	// create input and output pipes
	ir, iw, err := os.Pipe()
	if err != nil {
		return nil, ServerError{errPHPInit, err}
	}
	or, ow, err := os.Pipe()
	if err != nil {
		ir.Close()
		iw.Close()
		return nil, ServerError{errPHPInit, err}
	}

	// create a context to close the server
	context, cancel := context.WithCancel(context.Background())

	// start the shell process, which will close everything once done
	go func() {
		defer func() {
			ir.Close()
			iw.Close()
			or.Close()
			ow.Close()

			cancel()
		}()

		// start the server
		io := stream.NewIOStream(ow, nil, ir, 0)
		execPHP(io, serverPHP)
	}()

	// return a new server
	return &Server{
		in:  iw,
		out: or,
		c:   context,
	}, nil
}

// Server represents a server that executes code within a distillery.
// A typical use-case is to define functions using [MarshalEval], and then call those functions [MarshalCall].
//
// A nil Server will return [ErrServerBroken] on every function call.
type Server struct {
	m sync.Mutex

	in  io.WriteCloser
	out io.Reader
	c   context.Context
}

// MarshalEval evaluates code on the server and Marshals the result into value.
// When value is nil, the results are discarded.
//
// code is directly passed to php's "eval" function.
// as such any functions defined will remain in server memory.
//
// When an exception is thrown by the PHP Code, error is not nil, and dest remains unchanged.
func (server *Server) MarshalEval(value any, code string) error {
	server.m.Lock()
	defer server.m.Unlock()

	// quick hack: when the server is already done
	if err := server.c.Err(); err != nil {
		return errPHPClosed
	}

	// find a delimiter for the code, and then send
	delim := findDelimiter(code)
	io.WriteString(server.in, delim+"\n"+code+"\n"+delim+"\n")

	// read the next line (as a response)
	data, err := nobufio.ReadLine(server.out)
	if err != nil {
		return ServerError{Message: errPHPReceive, Err: err}
	}

	// read whatever we received
	var received [2]json.RawMessage
	if err := json.Unmarshal([]byte(data), &received); err != nil {
		return ServerError{Message: errPHPMarshal, Err: err}
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
func (server *Server) Eval(code string) (value any, err error) {
	err = server.MarshalEval(&value, code)
	return
}

// MarshalCall calls a previously defined function with the given arguments.
// Arguments are sent to php using json Marshal, and are 'json_decode'd on the php side.
//
// Return values are received as in [MarshalEval].
func (server *Server) MarshalCall(value any, function string, args ...any) error {

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
	return server.MarshalEval(value, code)
}

// Call is like [MarshalCall] but returns the return value of the function as an any
func (server *Server) Call(function string, args ...any) (value any, err error) {
	err = server.MarshalCall(&value, function, args...)
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
	server.m.Lock()
	defer server.m.Unlock()

	// if the context is already closed
	if err := server.c.Err(); err != nil {
		return errPHPClosed
	}

	server.in.Close()
	<-server.c.Done()

	return nil
}

//go:embed server.php
var serverPHP string

// pre-process the server.php code to make it shorter
func init() {
	// remove the first '<?php' line
	lines := strings.Split(serverPHP, "\n")[1:]
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	// remove comment lines
	lines = collection.Filter(lines, func(line string) bool {
		return !strings.HasPrefix(line, "//")
	})

	serverPHP = strings.Join(lines, "")
}
