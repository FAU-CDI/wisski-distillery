package instances

import (
	"github.com/tkw1536/goprogram/stream"
)

// Shell executes a shell command inside the instance.
func (wisski *WissKI) Shell(io stream.IOStream, argv ...string) (int, error) {
	return wisski.Barrel().Exec(io, "barrel", "/bin/sh", append([]string{"/user_shell.sh"}, argv...)...)
}
