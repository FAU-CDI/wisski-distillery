// Package execx defines extensions to the "os/exec" package
package execx

import (
	"os/exec"

	"github.com/tkw1536/goprogram/stream"
)

// ExecCommandError is returned by Exec when a command could not be executed.
// This typically hints that the executable cannot be found, but may have other causes.
const ExecCommandError = 127

// Exec executes a system command with the specified input/output streams, working directory, and arguments.
//
// If the command executes, it's exit code will be returned.
// If the command can not be executed, returns [ExecCommandError].
func Exec(io stream.IOStream, workdir string, exe string, argv ...string) int {
	// setup the command
	cmd := exec.Command(exe, argv...)
	cmd.Dir = workdir
	cmd.Stdin = io.Stdin
	cmd.Stdout = io.Stdout
	cmd.Stderr = io.Stderr

	// run it
	err := cmd.Run()

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

// MustExec is like Exec, except that it returns true if the command exited successfully, and else false.
func MustExec(io stream.IOStream, workdir string, exe string, argv ...string) bool {
	return Exec(io, workdir, exe, argv...) == 0
}
