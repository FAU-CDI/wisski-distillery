package execx

import (
	"github.com/tkw1536/goprogram/stream"
)

// Compose runs a docker-compose command in a specific directory, with the provided arguments and streams.
// It then waits for the process to exit, and returns the exit code.
func Compose(io stream.IOStream, workdir string, args ...string) int {
	return Exec(io, workdir, "docker", append([]string{"compose"}, args...)...)
}
