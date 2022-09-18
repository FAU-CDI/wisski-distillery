package component

import (
	"io/fs"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/unpack"
	"github.com/pkg/errors"
	"github.com/tkw1536/goprogram/stream"
)

// TODO: Move this package into components

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
func (is StackWithResources) Install(env environment.Environment, io stream.IOStream, context InstallationContext) error {
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
