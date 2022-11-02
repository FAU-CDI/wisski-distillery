package environment

import (
	"os/exec"

	"github.com/tkw1536/goprogram/stream"
)

// Exec executes a system command with the specified input/output streams, working directory, and arguments.
//
// If the command executes, it's exit code will be returned.
// If the command can not be executed, returns [ExecCommandError].
func (*Native) Exec(io stream.IOStream, workdir string, exe string, argv ...string) int {
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

func (n *Native) LookPathAbs(file string) (string, error) {
	path, err := exec.LookPath(file)
	if err != nil {
		return "", err
	}
	return n.Abs(path)
}
