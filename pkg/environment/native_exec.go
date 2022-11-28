package environment

import (
	"context"
	"os/exec"

	"github.com/FAU-CDI/wisski-distillery/pkg/cancel"
	"github.com/tkw1536/goprogram/stream"
)

// Exec executes a system command with the specified input/output streams, working directory, and arguments.
//
// If the command executes, it's exit code will be returned.
// If the command can not be executed, returns [ExecCommandError].
func (*Native) Exec(ctx context.Context, io stream.IOStream, workdir string, exe string, argv ...string) int {
	// setup the command
	cmd := exec.Command(exe, argv...)
	cmd.Dir = workdir
	cmd.Stdin = io.Stdin
	cmd.Stdout = io.Stdout
	cmd.Stderr = io.Stderr

	// run the process in a cancelable fashion
	err, cErr := cancel.WithContext(ctx, func(cancelable func()) error {
		// start the process
		err := cmd.Start()
		if err != nil {
			return err
		}

		// allow it to be cancellable
		cancelable()

		// and wait for the rest of the process
		return cmd.Wait()
	}, func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	})
	if err == nil {
		err = cErr
	}

	// non-zero exit
	if err, ok := err.(*exec.ExitError); ok {
		return err.ExitCode()
	}

	// unknown error
	if err != nil {
		return ExecCommandError
	}

	// everything is fine!
	return 0
}

func (n *Native) LookPathAbs(file string) (string, error) {
	path, err := exec.LookPath(file)
	if err != nil {
		return "", err
	}
	return n.Abs(path)
}
