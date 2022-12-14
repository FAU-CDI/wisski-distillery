package barrel

import (
	"context"

	"github.com/tkw1536/goprogram/stream"
)

// Shell executes a shell command inside the instance.
func (barrel *Barrel) Shell(ctx context.Context, io stream.IOStream, argv ...string) func() int {
	return barrel.Stack().Exec(ctx, io, "barrel", "/bin/sh", append([]string{"/user_shell.sh"}, argv...)...)
}
