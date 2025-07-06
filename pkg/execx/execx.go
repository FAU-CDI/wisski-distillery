// Package execx provides thin wrappers around the os.Exec package.
//
//spellchecker:words execx
package execx

// TODO: Move this to an external package

//spellchecker:words context errors exec path filepath time github wisski distillery internal wdlog pkglib stream
import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"go.tkw01536.de/pkglib/stream"
)

// CommandError is returned by Exec when a command could not be executed.
// This typically hints that the executable cannot be found, but may have other causes.
const CommandError = 127

// CommandErrorFunc always returns CommandError.
func CommandErrorFunc() int { return CommandError }

// Exec executes a system command with the specified input/output streams, working directory, and arguments.
// If the process is killed, and fails to exit, an internal amount of time is waited before it is closed.
//
// The command is started immediatly.
// The returned function is guaranteed to be non-nil and returns an exit code.
//
// If the command executes, the returns the exit code as soon as the process executes.
// If the command can not be executed, the returned function is [ExecCommandErrorFunc] and returns [ExecCommandError].
func Exec(ctx context.Context, io stream.IOStream, workdir string, exe string, argv ...string) func() int {
	// setup the command
	cmd := exec.CommandContext(ctx, exe, argv...)
	cmd.WaitDelay = time.Second
	cmd.Dir = workdir
	cmd.Stdin = io.Stdin
	cmd.Stdout = io.Stdout
	cmd.Stderr = io.Stderr

	// context is already cancelled, don't run it!
	if err := ctx.Err(); err != nil {
		return CommandErrorFunc
	}

	// start the command, but if something happens, return nil
	err := cmd.Start()
	wdlog.Of(ctx).Debug(
		"exec.Command.Start",
		"exe", exe,
		"argv", argv,
		"error", err,
	)
	if err != nil {
		return CommandErrorFunc
	}

	// create a new command
	return func() int {
		err := cmd.Wait()
		wdlog.Of(ctx).Debug(
			"exec.Command.Wait",
			"exe", exe,
			"argv", argv,
			"error", err,
		)

		// non-zero exit
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			return exit.ExitCode()
		}

		if err != nil {
			return CommandError
		}

		return 0
	}
}

// MustExec is like Exec, except that it returns true if the command exited successfully, and else false.
func MustExec(ctx context.Context, io stream.IOStream, workdir string, exe string, argv ...string) bool {
	return Exec(ctx, io, workdir, exe, argv...)() == 0
}

// LookPathAbs looks for the path named "file" and then resolves this path absolutely.
func LookPathAbs(file string) (string, error) {
	path, err := exec.LookPath(file)
	if err != nil {
		return "", fmt.Errorf("failed to find file: %w", err)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}
	return abs, nil
}
