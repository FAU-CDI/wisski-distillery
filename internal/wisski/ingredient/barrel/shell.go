package barrel

import "github.com/tkw1536/goprogram/stream"

// Shell executes a shell command inside the instance.
func (barrel *Barrel) Shell(io stream.IOStream, argv ...string) (int, error) {
	return barrel.Stack().Exec(io, "barrel", "/bin/sh", append([]string{"/user_shell.sh"}, argv...)...)
}
