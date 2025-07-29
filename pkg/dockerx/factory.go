// Package dockerx implements extended docker api functionality.
//
//spellchecker:words dockerx
package dockerx

//spellchecker:words pkglib errorsx
import (
	"fmt"

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

// WithFunc creates a new client from the factory, executes f, and closes the client.
func WithFunc[T any](factory Factory, f func(*Client) (T, error)) (t T, e error) {
	client, err := factory.NewClient()
	if err != nil {
		return t, fmt.Errorf("failed to create client: %w", err)
	}
	defer errorsx.Close(client, &e, "client")

	res, err := f(client)
	if err != nil {
		return t, fmt.Errorf("failed to execute func: %w", err)
	}
	return res, nil
}

// WithFunc0 execpt that func only returns an error.
func WithFunc0(factory Factory, f func(*Client) error) (e error) {
	client, err := factory.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer errorsx.Close(client, &e, "client")

	if err := f(client); err != nil {
		return fmt.Errorf("failed to execute func: %w", err)
	}
	return nil
}

func WithFunc2[T, U any](factory Factory, f func(*Client) (T, U, error)) (t T, u U, e error) {
	client, err := factory.NewClient()
	if err != nil {
		return t, u, fmt.Errorf("failed to create client: %w", err)
	}
	defer errorsx.Close(client, &e, "client")

	res1, res2, err := f(client)
	if err != nil {
		return t, u, fmt.Errorf("failed to execute func: %w", err)
	}
	return res1, res2, nil
}
