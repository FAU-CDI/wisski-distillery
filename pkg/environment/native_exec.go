package environment

import (
	"context"
	"os/exec"

	"github.com/rs/zerolog"
	"github.com/tkw1536/goprogram/stream"
)

// Exec executes a system command with the specified input/output streams, working directory, and arguments.
//
// The command is started immediatly.
// The returned function is guaranteed to be non-nil and returns an exit code.
//
// If the command executes, the returns the exit code as soon as the process executes.
// If the command can not be executed, the returned function is [ExecCommandErrorFunc] and returns [ExecCommandError].
func (*Native) Exec(ctx context.Context, io stream.IOStream, workdir string, exe string, argv ...string) func() int {
	// setup the command
	cmd := exec.Command(exe, argv...)
	cmd.Dir = workdir
	cmd.Stdin = io.Stdin
	cmd.Stdout = io.Stdout
	cmd.Stderr = io.Stderr

	// context is already cancelled, don't run it!
	if err := ctx.Err(); err != nil {
		return ExecCommandErrorFunc
	}

	// start the command, but if something happens, return nil
	err := cmd.Start()
	zerolog.Ctx(ctx).Debug().Str("exe", exe).Strs("argv", argv).Err(err).Msg("exec.Command.Start")
	if err != nil {
		return ExecCommandErrorFunc
	}

	waitdone := make(chan struct{}) // closed once Wait() below returns
	alldone := make(chan struct{})  // closed once the kill goroutine exits
	go func() {
		defer close(alldone)

		select {
		case <-ctx.Done():
			err := cmd.Process.Kill()
			zerolog.Ctx(ctx).Debug().Str("exe", exe).Strs("argv", argv).Err(err).Msg("exec.Command.Kill")
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
		zerolog.Ctx(ctx).Debug().Str("exe", exe).Strs("argv", argv).Err(err).Msg("exec.Command.Wait")

		// non-zero exit
		if err, ok := err.(*exec.ExitError); ok {
			return err.ExitCode()
		}

		if err != nil {
			return ExecCommandError
		}

		return 0
	}
}

func (n *Native) LookPathAbs(file string) (string, error) {
	path, err := exec.LookPath(file)
	if err != nil {
		return "", err
	}
	return n.Abs(path)
}
