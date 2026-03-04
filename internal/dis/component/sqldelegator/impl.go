package sqldelegator

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/stream"
	"go.tkw01536.de/pkglib/timex"
)

// Impl represents an sql implementation inside a docker container.
type Impl struct {
	Service   string
	OpenStack func() (*dockerx.Stack, error)

	QueryExecutable string
	DumpExecutable  string

	// Interval used to poll the sql database being available.
	PollInterval time.Duration

	// Timeout for waiting for the sql database to start.
	StartTimeout time.Duration
}

// NewImpl creates a new sql implementation.
func NewImpl(service string, openStack func() (*dockerx.Stack, error)) *Impl {
	return &Impl{
		Service:   service,
		OpenStack: openStack,

		QueryExecutable: "mariadb",
		DumpExecutable:  "mariadb-dump",

		PollInterval: 1 * time.Second,
		StartTimeout: 1 * time.Minute,
	}
}

// whileRunning executes the given function while the container is running, and the sql database itself is up.
//
// If the container is not running, it starts it, runs the function, and then stops it again once the function returns.
// If the container is already running, it simply waits for the sql database to be up, runs the function, and then returns.
// It checks if the sql database is up by executing a 'select 1' query.
func (sql *Impl) whileRunning(ctx context.Context, progress io.Writer, fn func(stack *dockerx.Stack) error) (e error) {
	// Open the stack
	stack, err := sql.OpenStack()
	if err != nil {
		return err
	}
	defer errorsx.Close(stack, &e, "stack")

	running, err := stack.Running(ctx, sql.Service)
	if err != nil {
		return fmt.Errorf("failed to check if container is running: %w", err)
	}

	// Start the containiner if it is not running
	if !running {
		if err := stack.Start(ctx, progress, sql.Service); err != nil {
			return fmt.Errorf("failed to start container: %w", err)
		}

		defer func() {
			if err := stack.Down(ctx, progress, sql.Service); err != nil {
				e = errorsx.Combine(e, fmt.Errorf("failed to stop container: %w", err))
			}
		}()
	}

	// Set a hard timeout for the wait operation.
	// This avoids having infinite loops.
	waitCtx, waitCtxCancel := context.WithTimeout(ctx, sql.StartTimeout)
	defer waitCtxCancel()

	// Wait for the sql database to be up
	connectErr := timex.TickUntilFunc(func(time.Time) bool {
		// We cannot use the Shell() function here, because that would recursively call this function.
		return stack.Exec(ctx, stream.NewIOStream(nil, progress, nil), dockerx.ExecOptions{
			Service: sql.Service,
			Cmd:     sql.QueryExecutable,
			Args:    []string{"-e", "select 1;"},
		})() == 0
	}, waitCtx, sql.PollInterval)
	if connectErr != nil {
		return fmt.Errorf("failed to wait for sql database to be up: %w", connectErr)
	}

	// and run the actual function
	return fn(stack)
}

// queries executes the given queries inside the sql implementation.
// If the container is not running, it is started automatically.
// They are run inside the default database, unless a different database is selected with a "USE database;" query.
func (sql *Impl) queries(ctx context.Context, progress io.Writer, queries ...string) error {
	return sql.whileRunning(ctx, progress, func(stack *dockerx.Stack) error {
		input := strings.NewReader(strings.Join(queries, ";\n"))
		io := stream.NewIOStream(progress, progress, input)

		code := stack.Exec(ctx, io, dockerx.ExecOptions{
			Service: sql.Service,
			Cmd:     sql.QueryExecutable,
		})()
		if code != 0 {
			return fmt.Errorf("failed to execute queries: %w", code)
		}
		return nil
	})
}
