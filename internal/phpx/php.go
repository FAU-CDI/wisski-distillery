// Package phpx provides functionalities for interacting with PHP code
//
//spellchecker:words phpx
package phpx

//spellchecker:words context pkglib stream
import (
	"context"

	"go.tkw01536.de/pkglib/stream"
)

// Executor represents anything that can spawn.
type Executor interface {
	// Spawn spawns a new (independent) process executing code.
	// It should return only once the execution terminates.
	Spawn(ctx context.Context, str stream.IOStream, code string) error
}

// SpawnFunc implements Executor.
type SpawnFunc func(ctx context.Context, str stream.IOStream, code string) error

func (sf SpawnFunc) Spawn(ctx context.Context, str stream.IOStream, code string) error {
	return sf(ctx, str, code)
}
