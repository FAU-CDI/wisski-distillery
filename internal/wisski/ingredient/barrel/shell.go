package barrel

import (
	"context"
	"fmt"

	"github.com/alessio/shellescape"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/stream"
)

type ExitError int

func (ee ExitError) Error() string {
	return fmt.Sprintf("Exited with code %d", int(ee))
}

func (ee ExitError) Code() exit.ExitCode {
	return exit.ExitCode(ee)
}

// Shell executes a shell with the given command line arguments inside the container.
// If an error occurs, it is of type ExitError.
func (barrel *Barrel) Shell(ctx context.Context, io stream.IOStream, argv ...string) error {
	code := barrel.Stack().Exec(ctx, io, "barrel", "/bin/sh", append([]string{"/user_shell.sh"}, argv...)...)()
	if code != 0 {
		return ExitError(code)
	}
	return nil
}

// ShellScript quotes the given command and executes it as a shell script inside the container.
func (barrel *Barrel) ShellScript(ctx context.Context, io stream.IOStream, commands ...string) error {
	command := shellescape.QuoteCommand(commands)
	return barrel.Shell(ctx, io, "-c", command)
}
