// Package stack implements a docker compose stack
package stack

import (
	"errors"

	"github.com/FAU-CDI/wisski-distillery/internal/execx"
	"github.com/tkw1536/goprogram/stream"
)

// Stack represents a 'docker compose' stack living in a specific directory
//
// NOTE(twiesing): In the current implementation this requires a 'docker' executable on the system.
// This executable must be capable of the 'docker compose' command.
// In the future the idea is to replace this with a native docker compose client.
type Stack struct {
	Name string // Name of this stack, TODO: Do we need this?
	Dir  string // Directory of this stack
}

var errStackUpdatePull = errors.New("Stack.Update: Pull returned non-zero exit code")
var errStackUpdateBuild = errors.New("Stack.Update: Build returned non-zero exit code")

// Update pulls, builds, and then optionally starts this stack.
// This does not have a direct 'docker compose' shell equivalent.
//
// See also Up.
func (ds Stack) Update(io stream.IOStream, start bool) error {
	if ds.compose(io, "pull") != 0 {
		return errStackUpdatePull
	}
	if ds.compose(io, "build", "--pull") != 0 {
		return errStackUpdateBuild
	}
	if start {
		return ds.Up(io)
	}
	return nil
}

var errStackUp = errors.New("Stack.Up: Up returned non-zero exit code")

// Up creates and starts the containers in this Stack.
// It is equivalent to 'docker compose up -d' on the shell.
func (ds Stack) Up(io stream.IOStream) error {
	if ds.compose(io, "up", "-d") != 0 {
		return errStackUp
	}
	return nil
}

// Exec executes an executable in the provided running service.
// It is equivalent to 'docker compose exec $service $executable $args...'.
//
// It returns the exit code of the process.
func (ds Stack) Exec(io stream.IOStream, service, executable string, args ...string) int {
	compose := []string{"exec"}
	if io.StdinIsATerminal() {
		compose = append(compose, "-ti")
	}
	compose = append(compose, executable)
	compose = append(compose, args...)
	return ds.compose(io, compose...)
}

// Run executes the provided service with the given executable.
// It is equivalent to 'docker compose run [--rm] $service $executable $args...'.
//
// It returns the exit code of the process.
func (ds Stack) Run(io stream.IOStream, autoRemove bool, service, command string, args ...string) int {
	compose := []string{"run"}
	if autoRemove {
		compose = append(compose, "--rm")
	}
	if !io.StdinIsATerminal() {
		compose = append(compose, "-T")
	}
	compose = append(compose, command)
	compose = append(compose, args...)
	return ds.compose(io, compose...)
}

var errStackRestart = errors.New("Stack.Restart: Restart returned non-zero exit code")

// Restart restarts all containers in this Stack.
// It is equivalent to 'docker compose restart' on the shell.
func (ds Stack) Restart(io stream.IOStream) error {
	if ds.compose(io, "restart") != 0 {
		return errStackRestart
	}
	return nil
}

var errStackDown = errors.New("Stack.Down: Down returned non-zero exit code")

// Down stops and removes all containers in this Stack.
// It is equivalent to 'docker compose down -v' on the shell.
func (ds Stack) Down(io stream.IOStream) error {
	if ds.compose(io, "down", "-v") != 0 {
		return errStackDown
	}
	return nil
}

// Compose executes a 'docker compose' command on this stack.
// TODO: This should be removed and replaced by an internal call directly to libcompose.
func (ds Stack) compose(io stream.IOStream, args ...string) int {
	// TODO: can we migrate to a built-in version of this?
	return execx.Compose(io, ds.Dir, args...)
}
