// Package phpx provides functionalities for interacting with PHP code
package phpx

import (
	"context"

	"github.com/tkw1536/pkglib/stream"
)

// Executor represents anything that can spawn
type Executor interface {
	// Spawn spawns a new (independent) process executing code.
	// It should return only once the execution terminates.
	Spawn(ctx context.Context, str stream.IOStream, code string) error
}

// SpawnFunc implements Executor
type SpawnFunc func(ctx context.Context, str stream.IOStream, code string) error

func (sf SpawnFunc) Spawn(ctx context.Context, str stream.IOStream, code string) error {
	return sf(ctx, str, code)
}
