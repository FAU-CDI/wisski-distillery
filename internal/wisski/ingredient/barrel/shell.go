//spellchecker:words barrel
package barrel

//spellchecker:words context essio shellescape github wisski distillery dockerx pkglib errorsx exit stream
import (
	"context"
	"fmt"

	"al.essio.dev/pkg/shellescape"
	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/stream"
)

type ExitError int

func (ee ExitError) Error() string {
	return fmt.Sprintf("command exited with code %d", int(ee))
}

func (ee ExitError) Code() exit.ExitCode {
	return exit.Code(int(ee))
}

// BashScript executes the given command as a bash script inside the container.
func (barrel *Barrel) BashScript(ctx context.Context, io stream.IOStream, commands ...string) (e error) {
	stack, err := barrel.OpenStack()
	if err != nil {
		return fmt.Errorf("failed to open stack: %w", err)
	}
	defer errorsx.Close(stack, &e, "stack")

	code := stack.Exec(
		ctx, io,
		dockerx.ExecOptions{
			Service: "barrel",
			User:    "www-data",

			Cmd:  "/bin/bash",
			Args: []string{"-c", shellescape.QuoteCommand(commands)},
		},
	)()
	if code != 0 {
		return ExitError(code)
	}
	return nil
}
