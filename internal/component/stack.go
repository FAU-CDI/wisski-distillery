// Package stack implements a docker compose stack
package component

import (
	"bufio"
	"bytes"
	"io/fs"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/unpack"
	"github.com/pkg/errors"
	"github.com/tkw1536/goprogram/stream"
)

// Stack represents a 'docker compose' stack living in a specific directory
//
// NOTE(twiesing): In the current implementation this requires a 'docker' executable on the system.
// This executable must be capable of the 'docker compose' command.
// In the future the idea is to replace this with a native docker compose client.
type Stack struct {
	Dir string // Directory this Stack is located in

	Env              environment.Environment
	DockerExecutable string // Path to the native docker executable to use
}

var errStackUpdatePull = errors.New("Stack.Update: Pull returned non-zero exit code")
var errStackUpdateBuild = errors.New("Stack.Update: Build returned non-zero exit code")

// Update pulls, builds, and then optionally starts this stack.
// This does not have a direct 'docker compose' shell equivalent.
//
// See also Up.
func (ds Stack) Update(io stream.IOStream, start bool) error {
	{
		code, err := ds.compose(io, "pull")
		if err != nil {
			return err
		}
		if code != 0 {
			return errStackUpdatePull
		}
	}

	{
		code, err := ds.compose(io, "build", "--pull")
		if err != nil {
			return err
		}
		if code != 0 {
			return errStackUpdateBuild
		}
	}
	if start {
		return ds.Up(io)
	}
	return nil
}

var errStackUp = errors.New("Stack.Up: Up returned non-zero exit code")

// Up creates and starts the containers in this Stack.
// It is equivalent to 'docker compose up --remove-orphans --detach' on the shell.
func (ds Stack) Up(io stream.IOStream) error {
	code, err := ds.compose(io, "up", "--remove-orphans", "--detach")
	if err != nil {
		return err
	}
	if code != 0 {
		return errStackUp
	}
	return nil
}

// Exec executes an executable in the provided running service.
// It is equivalent to 'docker compose exec $service $executable $args...'.
//
// It returns the exit code of the process.
func (ds Stack) Exec(io stream.IOStream, service, executable string, args ...string) (int, error) {
	compose := []string{"exec"}
	if io.StdinIsATerminal() {
		compose = append(compose, "-ti")
	}
	compose = append(compose, service)
	compose = append(compose, executable)
	compose = append(compose, args...)
	return ds.compose(io, compose...)
}

// Run executes the provided service with the given executable.
// It is equivalent to 'docker compose run [--rm] $service $executable $args...'.
//
// It returns the exit code of the process.
func (ds Stack) Run(io stream.IOStream, autoRemove bool, service, command string, args ...string) (int, error) {
	compose := []string{"run"}
	if autoRemove {
		compose = append(compose, "--rm")
	}
	if !io.StdinIsATerminal() {
		compose = append(compose, "-T")
	}
	compose = append(compose, service, command)
	compose = append(compose, args...)

	code, err := ds.compose(io, compose...)
	if err != nil {
		return environment.ExecCommandError, nil
	}
	return code, nil
}

var errStackRestart = errors.New("Stack.Restart: Restart returned non-zero exit code")

// Restart restarts all containers in this Stack.
// It is equivalent to 'docker compose restart' on the shell.
func (ds Stack) Restart(io stream.IOStream) error {
	code, err := ds.compose(io, "restart")
	if err != nil {
		return err
	}
	if code != 0 {
		return errStackRestart
	}
	return nil
}

var errStackPs = errors.New("Stack.Ps: Down returned non-zero exit code")

// Ps returns the ids of the containers currently running
func (ds Stack) Ps(io stream.IOStream) ([]string, error) {
	// create a buffer
	var buffer bytes.Buffer

	// read the ids from the command!
	code, err := ds.compose(io.Streams(&buffer, nil, nil, 0), "ps", "-q")
	if err != nil {
		return nil, err
	}
	if code != 0 {
		return nil, errStackPs
	}

	// scan each of the lines
	var results []string
	scanner := bufio.NewScanner(&buffer)
	for scanner.Scan() {
		if text := scanner.Text(); text != "" {
			results = append(results, text)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// return them!
	return results, nil
}

var errStackDown = errors.New("Stack.Down: Down returned non-zero exit code")

// Down stops and removes all containers in this Stack.
// It is equivalent to 'docker compose down -v' on the shell.
func (ds Stack) Down(io stream.IOStream) error {
	code, err := ds.compose(io, "down", "-v")
	if err != nil {
		return err
	}
	if code != 0 {
		return errStackDown
	}
	return nil
}

// compose executes a 'docker compose' command on this stack.
//
// NOTE(twiesing): Check if this can be replaced by an internal call to libcompose.
// But probably not.
func (ds Stack) compose(io stream.IOStream, args ...string) (int, error) {
	if ds.DockerExecutable == "" {
		var err error
		ds.DockerExecutable, err = ds.Env.LookPathAbs("docker")
		if err != nil {
			return environment.ExecCommandError, err
		}
	}
	return ds.Env.Exec(io, ds.Dir, ds.DockerExecutable, append([]string{"compose"}, args...)...), nil
}

// StackWithResources represents a Stack that can be automatically installed from a set of resources.
// See the [Install] method.
type StackWithResources struct {
	Stack

	// Installable enabled installing several resources from a (potentially embedded) filesystem.
	//
	// The Resources holds these, with appropriate resources specified below.
	// These all refer to paths within the Resource filesystem.
	Resources   fs.FS
	ContextPath string            // the 'docker compose' stack context, containing e.g. 'docker-compose.yml'.
	EnvPath     string            // the '.env' template, will be installed using [unpack.InstallTemplate].
	EnvContext  map[string]string // context when instantiating the '.env' template

	CopyContextFiles []string // Files to copy from the installation context

	MakeDirsPerm fs.FileMode // permission for diretories, defaults to [environment.DefaultDirCreate]
	MakeDirs     []string    // directories to ensure that exist

	TouchFiles []string // Files to 'touch', i.e. ensure that exist; guaranteed to be run after MakeDirs
}

// InstallationContext is a context to install data in
type InstallationContext map[string]string

// Install installs or updates this stack into the directory specified by stack.Stack().
//
// Installation is non-interactive, but will provide debugging output onto io.
// InstallationContext
func (is StackWithResources) Install(io stream.IOStream, context InstallationContext) error {
	env := is.Stack.Env
	if is.ContextPath != "" {
		// setup the base files
		if err := unpack.InstallDir(
			env,
			is.Dir,
			is.ContextPath,
			is.Resources,
			func(dst, src string) {
				io.Printf("[install] %s\n", dst)
			},
		); err != nil {
			return err
		}
	}

	// configure .env
	envDest := filepath.Join(is.Dir, ".env")
	if is.EnvPath != "" && is.EnvContext != nil {
		io.Printf("[config]  %s\n", envDest)
		if err := unpack.InstallTemplate(
			env,
			envDest,
			is.EnvContext,
			is.EnvPath,
			is.Resources,
		); err != nil {
			return err
		}
	}

	// make sure that certain dirs exist
	for _, name := range is.MakeDirs {
		// find the destination!
		dst := filepath.Join(is.Dir, name)

		io.Printf("[make]    %s\n", dst)
		if is.MakeDirsPerm == fs.FileMode(0) {
			is.MakeDirsPerm = environment.DefaultDirPerm
		}
		if err := env.MkdirAll(dst, is.MakeDirsPerm); err != nil {
			return err
		}
	}

	// copy files from the context!
	for _, name := range is.CopyContextFiles {
		// find the source!
		src, ok := context[name]
		if !ok {
			return errors.Errorf("Missing file from context: %q", src)
		}

		// find the destination!
		dst := filepath.Join(is.Dir, name)

		// copy over file from context
		io.Printf("[copy]    %s (from %s)\n", dst, src)
		if err := fsx.CopyFile(env, dst, src); err != nil {
			return errors.Wrapf(err, "Unable to copy file %s", src)
		}
	}

	// make sure that certain files exist
	for _, name := range is.TouchFiles {
		// find the destination!
		dst := filepath.Join(is.Dir, name)

		io.Printf("[touch]   %s\n", dst)
		if err := fsx.Touch(env, dst); err != nil {
			return err
		}
	}

	return nil
}
