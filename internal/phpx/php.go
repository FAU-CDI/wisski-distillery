// Package phpx provides functionalities for interacting with PHP code
package phpx

import "github.com/tkw1536/goprogram/stream"

// Executor represents anything that can spawn
type Executor interface {
	// Spawn spawns a new (independent) process executing code.
	// It should return only once the execution terminates.
	Spawn(str stream.IOStream, code string) error
}

// SpawnFunc implements Executor
type SpawnFunc func(str stream.IOStream, code string) error

func (sf SpawnFunc) Spawn(str stream.IOStream, code string) error {
	return sf(str, code)
}
