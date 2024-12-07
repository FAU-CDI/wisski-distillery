// Package stack implements a docker compose stack
//
//spellchecker:words component
package component

//spellchecker:words context path filepath github wisski distillery compose execx unpack errors pkglib umaskfree stream gopkg yaml
import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/pkg/compose"
	"github.com/FAU-CDI/wisski-distillery/pkg/execx"
	"github.com/FAU-CDI/wisski-distillery/pkg/unpack"
	"github.com/pkg/errors"
	"github.com/tkw1536/pkglib/fsx"
	"github.com/tkw1536/pkglib/fsx/umaskfree"
	"github.com/tkw1536/pkglib/stream"
	"gopkg.in/yaml.v3"
)

// Stack represents a 'docker compose' stack living in a specific directory
//
// NOTE(twiesing): In the current implementation this requires a 'docker' executable on the system.
// This executable must be capable of the 'docker compose' command.
// In the future the idea is to replace this with a native docker compose client.
type Stack struct {
	Dir string // Directory this Stack is located in

	DockerExecutable string // Path to the native docker executable to use
}

var errStackKill = errors.New("Stack.Kill: Kill returned non-zero exit code")

func (ds Stack) Kill(ctx context.Context, progress io.Writer, service string, signal os.Signal) error {
	code := ds.compose(ctx, stream.NonInteractive(progress), "kill", service, "-s", signal.String())()
	if code != 0 {
		return errStackKill
	}
	return nil
}

var errStackUpdatePull = errors.New("Stack.Update: Pull returned non-zero exit code")
var errStackUpdateBuild = errors.New("Stack.Update: Build returned non-zero exit code")

// Update pulls, builds, and then optionally starts this stack.
// This does not have a direct 'docker compose' shell equivalent.
//
// See also Up.
func (ds Stack) Update(ctx context.Context, progress io.Writer, start bool) error {
	if code := ds.compose(ctx, stream.NonInteractive(progress), "pull")(); code != 0 {
		return errStackUpdatePull
	}

	if code := ds.compose(ctx, stream.NonInteractive(progress), "build", "--pull")(); code != 0 {
		return errStackUpdateBuild
	}

	if start {
		return ds.Up(ctx, progress)
	}
	return nil
}

var errStackUp = errors.New("Stack.Up: Up returned non-zero exit code")

// Up creates and starts the containers in this Stack.
// It is equivalent to 'docker compose up --force-recreate --remove-orphans --detach' on the shell.
func (ds Stack) Up(ctx context.Context, progress io.Writer) error {
	if code := ds.compose(ctx, stream.NonInteractive(progress), "up", "--force-recreate", "--remove-orphans", "--detach")(); code != 0 {
		return errStackUp
	}
	return nil
}

// Exec executes an executable in the provided running service.
// It is equivalent to 'docker compose exec $service $executable $args...'.
//
// It returns the exit code of the process.
func (ds Stack) Exec(ctx context.Context, io stream.IOStream, service, executable string, args ...string) func() int {
	compose := []string{"exec"}
	if !io.StdinIsATerminal() {
		compose = append(compose, "-T")
	}

	compose = append(compose, service)
	compose = append(compose, executable)
	compose = append(compose, args...)

	return ds.compose(ctx, io, compose...)
}

type RunFlags struct {
	AutoRemove bool
	Detach     bool
}

// Run runs a command in a running container with the given executable.
// It is equivalent to 'docker compose run [--rm] $service $executable $args...'.
//
// It returns the exit code of the process.
func (ds Stack) Run(ctx context.Context, io stream.IOStream, flags RunFlags, service, command string, args ...string) (int, error) {
	compose := []string{"run"}
	if flags.AutoRemove {
		compose = append(compose, "--rm")
	}
	if !io.StdinIsATerminal() {
		compose = append(compose, "--no-TTY")
	}
	if flags.Detach {
		compose = append(compose, "--detach")
	}

	compose = append(compose, service, command)
	compose = append(compose, args...)

	code := ds.compose(ctx, io, compose...)()
	return code, nil
}

var errStackRestart = errors.New("Stack.Restart: Restart returned non-zero exit code")

// Restart restarts all containers in this Stack.
// It is equivalent to 'docker compose restart' on the shell.
func (ds Stack) Restart(ctx context.Context, progress io.Writer) error {
	code := ds.compose(ctx, stream.NonInteractive(progress), "restart")()
	if code != 0 {
		return errStackRestart
	}
	return nil
}

var errStackDown = errors.New("Stack.Down: Down returned non-zero exit code")

// Down stops and removes all containers in this Stack.
// It is equivalent to 'docker compose down -v' on the shell.
func (ds Stack) Down(ctx context.Context, progress io.Writer) error {
	code := ds.compose(ctx, stream.NonInteractive(progress), "down", "-v")()
	if code != 0 {
		return errStackDown
	}
	return nil
}

// DownAll stops and removes all containers in this Stack, and those not defined in the compose file.
// It is equivalent to 'docker compose down -v --remove-orphans' on the shell.
func (ds Stack) DownAll(ctx context.Context, progress io.Writer) error {
	code := ds.compose(ctx, stream.NonInteractive(progress), "down", "-v", "--remove-orphans")()
	if code != 0 {
		return errStackDown
	}
	return nil
}

// compose executes a 'docker compose' command on this stack.
//
// NOTE(twiesing): Check if this can be replaced by an internal call to libcompose.
// But probably not.
func (ds Stack) compose(ctx context.Context, io stream.IOStream, args ...string) func() int {
	if ds.DockerExecutable == "" {
		var err error
		ds.DockerExecutable, err = execx.LookPathAbs("docker")
		if err != nil {
			return execx.CommandErrorFunc
		}
	}
	return execx.Exec(ctx, io, ds.Dir, ds.DockerExecutable, append([]string{"compose"}, args...)...)
}

// StackWithResources represents a Stack that can be automatically installed from a set of resources.
// See the [Install] method.
type StackWithResources struct {
	Stack

	// Installable enabled installing several resources from a (potentially embedded) filesystem.
	//
	// The Resources holds these, with appropriate resources specified below.
	// These all refer to paths within the Resource filesystem.
	Resources fs.FS

	ContextPath string // the 'docker compose' stack context. May or may not contain 'docker-compose.yml'

	ComposerYML func(*yaml.Node) (*yaml.Node, error) // update 'docker-compose.yml', if no 'docker-compose.yml' exists, the passed node is nil.

	EnvContext map[string]string // context when instantiating the '.env' template

	CopyContextFiles []string // Files to copy from the installation context

	MakeDirsPerm fs.FileMode // permission for dirctories, defaults to [environment.DefaultDirCreate]
	MakeDirs     []string    // directories to ensure that exist

	TouchFilesPerm fs.FileMode       // permission for new files to touch or create, defaults to [environment.DefaultFileCreate]
	TouchFiles     []string          // Files to 'touch', i.e. ensure that exist; guaranteed to be run after MakeDirs
	CreateFiles    map[string]string // Files to 'create' but not update after they are setup; guaranteed to be run after MakeDirs
}

// InstallationContext is a context to install data in
type InstallationContext map[string]string

// Install installs or updates this stack into the directory specified by stack.Stack().
//
// Installation is non-interactive, but will provide debugging output onto io.
// InstallationContext
func (is StackWithResources) Install(ctx context.Context, progress io.Writer, context InstallationContext) error {
	if is.ContextPath != "" {
		// setup the base files
		if err := unpack.InstallDir(
			is.Dir,
			is.ContextPath,
			is.Resources,
			func(dst, src string) {
				fmt.Fprintf(progress, "[install] %s\n", dst)
			},
		); err != nil {
			return err
		}
	} else {
		if err := umaskfree.MkdirAll(is.Dir, umaskfree.DefaultDirPerm); err != nil {
			return err
		}
	}

	dockerComposeYML := filepath.Join(is.Dir, "docker-compose.yml")

	// update the docker compose file
	if is.ComposerYML != nil {
		fmt.Fprintf(progress, "[install] %s\n", dockerComposeYML)
		if err := doComposeFile(dockerComposeYML, is.ComposerYML); err != nil {
			return err
		}
	}

	if err := addComposeFileHeader(dockerComposeYML); err != nil {
		fmt.Fprintf(progress, "[update] %s\n", dockerComposeYML)
		return err
	}

	// configure .env
	envDest := filepath.Join(is.Dir, ".env")
	if is.EnvContext != nil {
		fmt.Fprintf(progress, "[config]  %s\n", envDest)

		if err := writeEnvFile(envDest, is.TouchFilesPerm, is.EnvContext); err != nil {
			return err
		}
	}

	// make sure that certain dirs exist
	for _, name := range is.MakeDirs {
		// find the destination!
		dst := filepath.Join(is.Dir, name)

		fmt.Fprintf(progress, "[make]    %s\n", dst)
		if is.MakeDirsPerm == fs.FileMode(0) {
			is.MakeDirsPerm = umaskfree.DefaultDirPerm
		}
		if err := umaskfree.MkdirAll(dst, is.MakeDirsPerm); err != nil {
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
		fmt.Fprintf(progress, "[copy]    %s (from %s)\n", dst, src)
		if err := umaskfree.CopyFile(ctx, dst, src); err != nil {
			return errors.Wrapf(err, "Unable to copy file %s", src)
		}
	}

	// touch files that should be created empty
	for _, name := range is.TouchFiles {
		// find the destination!
		dst := filepath.Join(is.Dir, name)

		fmt.Fprintf(progress, "[touch]   %s\n", dst)
		if err := umaskfree.Touch(dst, umaskfree.DefaultFilePerm); err != nil {
			return err
		}
	}
	// make sure that certain files exist
	for name, content := range is.CreateFiles {
		// find the destination!
		dst := filepath.Join(is.Dir, name)

		exists, err := fsx.Exists(dst)
		if err != nil {
			return err
		}

		// create the file if it doesn't exist
		if !exists {
			fmt.Fprintf(progress, "[create]   %s\n", dst)
			if err := umaskfree.WriteFile(dst, []byte(content), umaskfree.DefaultFilePerm); err != nil {
				return err
			}
		} else {
			fmt.Fprintf(progress, "[skip]   %s\n", dst)
		}
	}

	// check that the stack can be loaded
	{
		fmt.Fprintln(progress, "[checking]")
		_, err := compose.Open(is.Dir)
		if err != nil {
			return err
		}
	}

	return nil
}

const composeFileHeader = "# This file was automatically created and is updated by the distillery; DO NOT EDIT.\n\n"

// addComposeFileHeader adds a header to the 'docker-compose.yml' file
// indicating it is automatically created
func addComposeFileHeader(path string) error {
	// read existing bytes
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// overwrite the file
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, umaskfree.DefaultFilePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	// write the header
	if _, err := f.WriteString(composeFileHeader); err != nil {
		return nil
	}

	// write the original content
	if _, err := f.Write(bytes); err != nil {
		return err
	}

	return nil
}

// doComposeFile updates the compose file using the update function.
//
// If the file at path already exists, calls update with the existing data.
// if the file at path does not exist, calls update(nil).
func doComposeFile(path string, update func(node *yaml.Node) (*yaml.Node, error)) error {
	var mode fs.FileMode
	var node *yaml.Node

	{
		stat, err := os.Stat(path)
		switch {
		case err == nil:
			// file exists => use the previous file mode
			mode = stat.Mode()

			// read the yaml bytes
			bytes, err := os.ReadFile(path)
			if err != nil {
				return errors.Wrap(err, "unable to read existing file")
			}

			// unmarshal it into a node, or bail out!
			node = new(yaml.Node)
			if err := yaml.Unmarshal(bytes, node); err != nil {
				return errors.Wrap(err, "unable to unmarshal existing file")
			}
		case errors.Is(err, fs.ErrNotExist):
			// file does not exist => use default mode
			mode = umaskfree.DefaultFilePerm

			// use a nil existing node
			node = nil
		default:
			return err
		}
	}

	// update the node
	node, err := update(node)
	if err != nil {
		return errors.Wrap(err, "update function failed")
	}

	// re-encode the bytes
	result, err := yaml.Marshal(node)
	if err != nil {
		return errors.Wrap(err, "failed to re-marshal")
	}

	// write the bytes back!
	return umaskfree.WriteFile(path, result, mode)
}

// writeEnvFile writes an environment file
func writeEnvFile(path string, perm fs.FileMode, variables map[string]string) error {
	// create the environment file
	file, err := umaskfree.Create(path, perm)
	if err != nil {
		return err
	}
	defer file.Close()

	// write the file!
	_, err = compose.WriteEnvFile(file, variables)
	if err != nil {
		return err
	}

	// and return nil
	return nil
}
