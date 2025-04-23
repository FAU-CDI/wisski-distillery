// Package execx provides thin wrappers around the os.Exec package.
//
//spellchecker:words execx
package execx

// TODO: Move this to an external package

//spellchecker:words context exec path filepath github wisski distillery internal wdlog pkglib stream
import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/tkw1536/pkglib/stream"
)

// CommandError is returned by Exec when a command could not be executed.
// This typically hints that the executable cannot be found, but may have other causes.
const CommandError = 127

// CommandErrorFunc always returns CommandError.
func CommandErrorFunc() int { return CommandError }

// Exec executes a system command with the specified input/output streams, working directory, and arguments.
//
// The command is started immediatly.
// The returned function is guaranteed to be non-nil and returns an exit code.
//
// If the command executes, the returns the exit code as soon as the process executes.
// If the command can not be executed, the returned function is [ExecCommandErrorFunc] and returns [ExecCommandError].
func Exec(ctx context.Context, io stream.IOStream, workdir string, exe string, argv ...string) func() int {
	// setup the command
	cmd := exec.Command(exe, argv...)
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

	waitdone := make(chan struct{}) // closed once Wait() below returns
	alldone := make(chan struct{})  // closed once the kill goroutine exits
	go func() {
		defer close(alldone)

		select {
		case <-ctx.Done():
			err := cmd.Process.Kill()
			wdlog.Of(ctx).Debug(
				"exec.Command.Kill",
				"exe", exe,
				"argv", argv,
				"error", err,
			)
		case <-waitdone:
		}
	}()

	// create a new command
	return func() int {
		defer func() {
			// wait for the goroutine to exit
			close(waitdone)
			<-alldone
		}()

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
