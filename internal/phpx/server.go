//spellchecker:words phpx
package phpx

//spellchecker:words bytes compress flate context encoding base json regexp slices strings sync embed github pkglib errorsx lazy status stream
import (
	"bytes"
	"compress/flate"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strings"
	"sync"

	_ "embed"

	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/lazy"
	"go.tkw01536.de/pkglib/status"
	"go.tkw01536.de/pkglib/stream"
)

// Server represents a server that executes PHP code.
// A typical use-case is to define functions using [MarshalEval], and then call those functions [MarshalCall].
//
// A server, once used, should be closed using the [Close] method.
//
//nolint:containedctx
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
	in    io.WriteCloser // input sent to the server
	lines chan string

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

		// create a context to close the server
		context, cancel := context.WithCancel(server.Context)
		server.cancel = cancel

		// store server props
		server.in = iw
		server.c = context

		server.lines = make(chan string, 1)
		lb := status.LineBuffer{
			Line: func(line string) {
				select {
				case server.lines <- line:
				default:
				}
			},
			CloseLine: func() {
				close(server.lines)
			},
			FlushLineOnClose: true,
		}

		// start the shell process, which will close everything once done
		go func() {
			defer func() {
				// TODO: is there a reasonable way to report this error?
				// via the logger perhaps?
				_ = ir.Close()
				_ = iw.Close()
				_ = lb.Close()

				server.cancel()
			}()

			// start the actual server
			io := stream.NewIOStream(&lb, nil, ir)
			err := server.Executor.Spawn(server.c, io, serverPHP)
			server.err.Set(ServerError{Message: errClosed, Err: err})
		}()
	})

	return server.err.Get(nil) //nolint:wrapcheck
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

	server.m.Lock()
	defer server.m.Unlock()

	// when the server is already done
	if err := server.c.Err(); err != nil {
		return ServerError{Message: errClosed}
	}

	// encode a message to the server!
	if err := server.encode(server.in, code); err != nil {
		server.cancel()
		return ServerError{Message: errSend, Err: err}
	}

	var data string
	var ok bool

	// read the next line from the server script
	select {
	case data, ok = <-server.lines:
	case <-server.c.Done():
	}

	if !ok {
		return ServerError{Message: errReceive, Err: io.EOF}
	}

	// decode the response
	var received [2]json.RawMessage
	if err := server.decode(&received, []byte(data)); err != nil {
		return ServerError{Message: errReceive, Err: err}
	}

	// check if there was an error
	var errString string
	if err := json.Unmarshal(received[1], &errString); err == nil && errString != "" {
		return ThrowableError(errString)
	}

	// special case: no return value => no unmarshaling needed
	if value == nil {
		return nil
	}

	// read the actual result!
	err := json.Unmarshal(received[0], value)
	if err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}
	return nil
}

// - decode json (opposite of php's "json_encode").
func (*Server) decode(dest *[2]json.RawMessage, message []byte) (e error) {
	// decode base64
	raw := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(message))

	// unpack gzip
	unpacker := flate.NewReader(raw)
	defer errorsx.Close(unpacker, &e, "unpacker")

	// and read the value
	decoder := json.NewDecoder(unpacker)
	if err := decoder.Decode(dest); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}
	return nil
}

// - encode base64 (opposite of php's "base64_decode").
func (*Server) encode(dest io.WriteCloser, code string) (e error) {
	// write a final newline at the end!
	defer func() {
		if e != nil {
			return
		}
		_, e = dest.Write([]byte("\n"))
	}()

	// base64 encode all the things!
	encoder := base64.NewEncoder(base64.StdEncoding, dest)
	defer errorsx.Close(encoder, &e, "encoder")

	// compress all the things!
	compressor, err := flate.NewWriter(encoder, 9)
	if err != nil {
		return fmt.Errorf("failed to create compressor: %w", err)
	}
	defer errorsx.Close(compressor, &e, "compressor")

	// do the write!
	_, e = compressor.Write([]byte(code))
	if e != nil {
		e = fmt.Errorf("failed to write to compressor: %w", e)
	}
	return e
}

// Eval is like [MarshalEval], but returns the value as an any.
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

// Call is like [MarshalCall] but returns the return value of the function as an any.
func (server *Server) Call(ctx context.Context, function string, args ...any) (value any, err error) {
	err = server.MarshalCall(ctx, &value, function, args...)
	return
}

// Close closes this server and prevents any further code from being run.
func (server *Server) Close() error {
	if err := server.prepare(); err != nil {
		return fmt.Errorf("failed to prepeare server: %w", err)
	}

	server.m.Lock()
	defer server.m.Unlock()

	// if the context is already closed
	if err := server.c.Err(); err != nil {
		return ServerError{Message: errClosed}
	}

	err := server.in.Close()
	<-server.c.Done()

	if err != nil {
		return fmt.Errorf("suspicous close of server input: %w", err)
	}
	return nil
}

//go:embed server.php
var serverPHP string

// pre-process the server.php code to make it shorter.
func init() {
	minifier := regexp.MustCompile(`\s*([=)(.,{}])\s*`)

	// remove the first '<?php' line
	lines := strings.Split(serverPHP, "\n")[1:]
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	// remove comment lines
	lines = slices.DeleteFunc(lines, func(line string) bool {
		return strings.HasPrefix(line, "//")
	})

	serverPHP = minifier.ReplaceAllString(strings.Join(lines, ""), "$1")
}
