// Package dockerx implements extended docker api functionality.
//
//spellchecker:words dockerx
package dockerx

//spellchecker:words pkglib errorsx
import (
	"context"
	"fmt"
	"io"
	"strings"

	"go.tkw01536.de/pkglib/errorsx"
)

// Factory can create docker clients.
type Factory interface {
	NewClient() (*Client, error)
}

// NewStack creates a stack from a factory and a directory.
func NewStack(factory Factory, dir string) (*Stack, error) {
	client, err := factory.NewClient()
	if err != nil {
		return nil, fmt.Errorf("factory failed to create client: %w", err)
	}
	return &Stack{
		Dir:    dir,
		Client: client,
	}, nil
}

// Do opens a stack using the given function, ensures that all of the services are running, and then runs the given function.
// If the services were not running, it starts them and then stops them at the end.
func Do(ctx context.Context, progress io.Writer, allowCreate bool, open func() (*Stack, error), f func(*Stack) error, services ...string) (e error) {
	stack, err := open()
	if err != nil {
		return fmt.Errorf("failed to open stack: %w", err)
	}
	defer errorsx.Close(stack, &e, "stack")

	// find services that aren't running
	toStart := make([]string, 0, len(services))
	for _, service := range services {
		running, err := stack.Running(ctx, service)
		if err != nil {
			return fmt.Errorf("failed to check if service %q is running: %w", service, err)
		}
		if !running {
			toStart = append(toStart, service)
		}
	}

	// if we have services to start, start them and then stop them again at the end!
	if len(toStart) > 0 {
		if !allowCreate {
			return fmt.Errorf("services %s are not running", strings.Join(toStart, ", "))
		}

		if err := stack.Start(ctx, progress, toStart...); err != nil {
			return fmt.Errorf("failed to start services %s: %w", strings.Join(toStart, ", "), err)
		}

		defer func() {
			if err := stack.Down(ctx, progress, toStart...); err != nil {
				e = errorsx.Combine(e, fmt.Errorf("failed to stop services %s: %w", strings.Join(toStart, ", "), err))
			}
		}()
	}

	// and do the actual running!
	if err := f(stack); err != nil {
		return fmt.Errorf("failed to run: %w", err)
	}
	return nil
}
