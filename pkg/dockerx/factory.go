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
