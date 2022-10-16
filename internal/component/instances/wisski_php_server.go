package instances

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	_ "embed"

	"github.com/alessio/shellescape"
	"github.com/tkw1536/goprogram/lib/collection"
	"github.com/tkw1536/goprogram/lib/nobufio"
	"github.com/tkw1536/goprogram/stream"
)

// Common PHP Error
var (
	errPHPInit    = "Unable to initialize"
	errPHPMarshal = "Marshal failed"
	errPHPInvalid = PHPServerError{Message: "Invalid code to execute"}
	errPHPReceive = "Failed to receive response"
	errPHPClosed  = PHPServerError{Message: "Server closed"}
)

// PHPError represents an error during PHPServer logic
type PHPServerError struct {
	Message string
	Err     error
}

func (err PHPServerError) Unwrap() error {
	return err.Err
}

func (err PHPServerError) Error() string {
	if err.Err == nil {
		return fmt.Sprintf("PHPServer: %s", err.Message)
	}
	return fmt.Sprintf("PHPServer: %s: %s", err.Message, err.Err)
}

// PHPThrowable represents an error during php code
type PHPThrowable string

func (throwable PHPThrowable) Error() string {
	return string(throwable)
}

// NewPHPServer returns a new server that can execute code within this distillery.
// When err == nil, the caller must call server.Close().
//
// See [PHPServer].
func (wisski *WissKI) NewPHPServer() (*PHPServer, error) {
	// create input and output pipes
	ir, iw, err := os.Pipe()
	if err != nil {
		return nil, PHPServerError{errPHPInit, err}
	}
	or, ow, err := os.Pipe()
	if err != nil {
		ir.Close()
		iw.Close()
		return nil, PHPServerError{errPHPInit, err}
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
		wisski.Shell(io, "-c", shellescape.QuoteCommand([]string{"drush", "php:eval", serverPHP}))
	}()

	// return the seerver
	return &PHPServer{
		in:  iw,
		out: or,
		c:   context,
	}, nil
}

// PHPServer represents a server that executes code within a distillery.
// A typical use-case is to define functions using [MarshalEval], and then call those functions [MarshalCall].
//
// A nil PHPServer will return [ErrServerBroken] on every function call.
type PHPServer struct {
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
func (server *PHPServer) MarshalEval(value any, code string) error {
	server.m.Lock()
	defer server.m.Unlock()

	// quick hack: when the server is already done
	if err := server.c.Err(); err != nil {
		return errPHPClosed
	}

	// marshal the code, and send it to the server
	bytes, err := json.Marshal(code)
	if err != nil {
		return PHPServerError{Message: errPHPMarshal, Err: err}
	}

	// send it to the server
	io.WriteString(server.in, string(bytes)+"\n")

	// read the next line (as a response)
	data, err := nobufio.ReadLine(server.out)
	if err != nil {
		return PHPServerError{Message: errPHPReceive, Err: err}
	}

	// read whatever we received
	var received [2]json.RawMessage
	if err := json.Unmarshal([]byte(data), &received); err != nil {
		return PHPServerError{Message: errPHPMarshal, Err: err}
	}

	// check if there was an error
	var errString string
	if err := json.Unmarshal(received[1], &errString); err == nil && errString != "" {
		return PHPThrowable(errString)
	}

	// special case: no return value => no unmarshaling needed
	if value == nil {
		return nil
	}

	// read the actual result!
	return json.Unmarshal(received[0], value)
}

// Eval is like [MarshalEval], but returns the value as an any
func (server *PHPServer) Eval(code string) (value any, err error) {
	err = server.MarshalEval(&value, code)
	return
}

// MarshalCall calls a previously defined function with the given arguments.
// Arguments are sent to php using json Marshal, and are 'json_decode'd on the php side.
//
// Return values are received as in [MarshalEval].
func (server *PHPServer) MarshalCall(value any, function string, args ...any) error {
	// marshal a code for the call
	userFunction, err := marshalPHP(function)
	if err != nil {
		return PHPServerError{Message: errPHPMarshal, Err: err}
	}
	userFunctionArgs, err := marshalPHP(args)
	if err != nil {
		return PHPServerError{Message: errPHPMarshal, Err: err}
	}
	code := "return call_user_func_array(" + userFunction + "," + userFunctionArgs + ");"

	// and return the evaluated code!
	return server.MarshalEval(value, code)
}

// Call is like [MarshalCall] but returns the return value of the function as an any
func (server *PHPServer) Call(function string, args ...any) (value any, err error) {
	err = server.MarshalCall(&value, function, args...)
	return
}

const marshalRune = 'F' // press to pay respect

// marshalPHP marshals some data which can be marshaled using [json.Encode] into a PHP Expression.
// the string can be safely used directly within php.
func marshalPHP(data any) (string, error) {
	// this function uses json as a data format to transport the data into php.
	// then we build a heredoc to encode it safely, and decode it in php

	// Step 1: Encode the data as json
	jbytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	jstring := string(jbytes)

	// Step 2: Find a delimiter for the heredoc.
	// Step 2a: Find the longest sequence of [marshalRune]s inside the encoded string.
	var current, longest int
	for _, r := range jstring {

		if r == marshalRune {
			current++
		} else {
			current = 0
		}

		if current > longest {
			longest = current
		}
	}
	// Step 2b: Build a string of marshalRune that is one longer!
	delim := strings.Repeat(string(marshalRune), longest+1)

	// Step 3: Assemble the encoded string!
	result := "call_user_func(function(){$x=<<<'" + delim + "'\n" + jstring + "\n" + delim + ";return json_decode(trim($x));})" // press to doubt
	return result, nil
}

// Close closes this server and prevents any further code from being run.
func (server *PHPServer) Close() error {
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

// ExecPHPScript executes the PHP code as a script on the given server.
// When server is nil, creates a new server and automatically closes it after execution.
//
// The script should define a function called entrypoint, and may define additional functions.
//
// Code must start with "<?php" and may not contain a closing tag.
// Code is expected not to mess with PHPs output buffer.
// Code should not contain user input.
// Code breaking these conventions may or may not result in an error.
//
// It's arguments are encoded as json using [json.Marshal] and decoded within php.
//
// The return value of the function is again marshaled with json and returned to the caller.
//
// Calling this function is inefficient, and a [NewPHPServer] call should be prefered instead.
func (wisski *WissKI) ExecPHPScript(server *PHPServer, value any, code string, entrypoint string, args ...any) (err error) {
	if server == nil {
		server, err = wisski.NewPHPServer()
		if err != nil {
			return
		}
		defer server.Close()
	}

	if err := server.MarshalEval(nil, strings.TrimPrefix(code, "<?php")); err != nil {
		return err
	}

	return server.MarshalCall(value, entrypoint, args...)
}

//go:embed php/server.php
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
