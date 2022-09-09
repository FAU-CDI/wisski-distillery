package stack

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/fsx"
	"github.com/pkg/errors"
	"github.com/tkw1536/goprogram/stream"
)

// Installable represents a Stack that can be automatically installed from a set of resources
// See the [Install] method.
type Installable struct {
	Stack

	ContextResource string // Path to the resource containing 'docker compose' context

	EnvFileResource string            // Path to the resource containing dynamically generated env file
	EnvFileContext  map[string]string // Context of variables to replace in the env file

	CopyContextFiles []string // Files to copy from the installation context

	MakeDirsPerm fs.FileMode // permission for diretories, defaults to fs.ModeDir
	MakeDirs     []string    // directories to ensure that exist

	TouchFiles []string // Files to 'touch', i.e. ensure that exist; guaranteed to be run after MakeDirs
}

// InstallationContext is a context to install data in
type InstallationContext map[string]string

// Install installs or updates this stack into the directory specified by stack.Stack().
//
// Installation is non-interactive, but will provide debugging output onto io.
// InstallationContext
func (is Installable) Install(io stream.IOStream, context InstallationContext) error {
	// setup the base files
	if err := distillery.InstallResource(
		is.Dir,
		is.ContextResource,
		func(dst, src string) {
			io.Printf("[install] %s\n", dst)
		},
	); err != nil {
		return err
	}

	// configure .env
	envDest := filepath.Join(is.Dir, ".env")
	if is.EnvFileResource != "" && is.EnvFileContext != nil {
		io.Printf("[config]  %s\n", envDest)
		if err := distillery.InstallTemplate(
			envDest,
			is.EnvFileResource,
			is.EnvFileContext,
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
			is.MakeDirsPerm = fs.ModeDir
		}
		if err := os.MkdirAll(dst, is.MakeDirsPerm); err != nil {
			return err
		}
	}

	// copy files from the context!
	for _, name := range is.CopyContextFiles {
		// find the source!
		src, ok := context[name]
		if !ok {
			return errors.Errorf("Missing file from context: %s", src)
		}

		// find the destination!
		dst := filepath.Join(is.Dir, name)

		// copy over file from context
		io.Printf("[copy]    %s (from %s)\n", dst, src)
		if err := fsx.CopyFile(dst, src); err != nil {
			return errors.Wrapf(err, "Unable to copy file %s", src)
		}
	}

	// make sure that certain files exist
	for _, name := range is.TouchFiles {
		// find the destination!
		dst := filepath.Join(is.Dir, name)

		io.Printf("[touch]   %s\n", dst)
		if err := fsx.Touch(dst); err != nil {
			return err
		}
	}

	return nil
}
