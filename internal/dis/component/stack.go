// Package stack implements a docker compose stack
//
//spellchecker:words component
package component

//spellchecker:words context errors path filepath github wisski distillery internal dockerenv dockerx unpack pkglib errorsx umaskfree gopkg yaml
import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/dockerenv"
	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
	"github.com/FAU-CDI/wisski-distillery/pkg/unpack"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/fsx"
	"go.tkw01536.de/pkglib/fsx/umaskfree"
	"gopkg.in/yaml.v3"
)

// StackWithResources represents a Stack that can be automatically installed from a set of resources.
// See the [Install] method.
type StackWithResources struct {
	*dockerx.Stack

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

// InstallationContext is a context to install data in.
type InstallationContext map[string]string

type fileMissingFromContextError string

func (fem fileMissingFromContextError) Error() string {
	return fmt.Sprintf("file missing from context: %q", string(fem))
}

// InstallationContext.
func (is StackWithResources) Install(ctx context.Context, progress io.Writer, context InstallationContext) error {
	if is.ContextPath != "" {
		// setup the base files
		if err := unpack.InstallDir(
			is.Dir,
			is.ContextPath,
			is.Resources,
			func(dst, src string) {
				// #nosec G103
				fmt.Fprintf(progress, "[install] %s\n", dst) //nolint:errcheck // no way to report error
			},
		); err != nil {
			return fmt.Errorf("failed to install directory: %w", err)
		}
	} else {
		if err := umaskfree.MkdirAll(is.Dir, umaskfree.DefaultDirPerm); err != nil {
			return fmt.Errorf("failed to create installation directory: %w", err)
		}
	}

	dockerComposeYML := filepath.Join(is.Dir, "docker-compose.yml")

	// update the docker compose file
	if is.ComposerYML != nil {
		if _, err := fmt.Fprintf(progress, "[install] %s\n", dockerComposeYML); err != nil {
			return fmt.Errorf("failed to log progress: %w", err)
		}

		if err := doComposeFile(dockerComposeYML, is.ComposerYML); err != nil {
			return fmt.Errorf("failed to update compose file: %w", err)
		}
	}

	if err := addComposeFileHeader(dockerComposeYML); err != nil {
		err = fmt.Errorf("failed to update docker compose yml: %w", err)
		if _, err2 := fmt.Fprintf(progress, "[update] %s\n", dockerComposeYML); err2 != nil {
			err = errorsx.Combine(
				err,
				fmt.Errorf("failed to log progress: %w", err2),
			)
		}
		return err
	}

	// configure .env
	envDest := filepath.Join(is.Dir, ".env")
	if is.EnvContext != nil {
		if _, err := fmt.Fprintf(progress, "[config]  %s\n", envDest); err != nil {
			return fmt.Errorf("failed to log progress: %w", err)
		}

		if err := writeEnvFile(envDest, is.TouchFilesPerm, is.EnvContext); err != nil {
			return fmt.Errorf("failed to write environment file: %w", err)
		}
	}

	// make sure that certain dirs exist
	for _, name := range is.MakeDirs {
		// find the destination!
		dst := filepath.Join(is.Dir, name)

		if _, err := fmt.Fprintf(progress, "[make]    %s\n", dst); err != nil {
			return fmt.Errorf("failed to log progress: %w", err)
		}
		if is.MakeDirsPerm == fs.FileMode(0) {
			is.MakeDirsPerm = umaskfree.DefaultDirPerm
		}
		if err := umaskfree.MkdirAll(dst, is.MakeDirsPerm); err != nil {
			return fmt.Errorf("failed to create directory %q: %w", dst, err)
		}
	}

	// copy files from the context!
	for _, name := range is.CopyContextFiles {
		// find the source!
		src, ok := context[name]
		if !ok {
			return fileMissingFromContextError(src)
		}

		// find the destination!
		dst := filepath.Join(is.Dir, name)

		// copy over file from context
		if _, err := fmt.Fprintf(progress, "[copy]    %s (from %s)\n", dst, src); err != nil {
			return fmt.Errorf("failed to report progress: %w", err)
		}
		if err := umaskfree.CopyFile(ctx, dst, src); err != nil {
			return fmt.Errorf("unable to copy file %s: %w", src, err)
		}
	}

	// touch files that should be created empty
	for _, name := range is.TouchFiles {
		// find the destination!
		dst := filepath.Join(is.Dir, name)

		if _, err := fmt.Fprintf(progress, "[touch]   %s\n", dst); err != nil {
			return fmt.Errorf("failed to report progress: %w", err)
		}
		if err := umaskfree.Touch(dst, umaskfree.DefaultFilePerm); err != nil {
			return fmt.Errorf("failed to touch %q: %w", dst, err)
		}
	}
	// make sure that certain files exist
	for name, content := range is.CreateFiles {
		// find the destination!
		dst := filepath.Join(is.Dir, name)

		exists, err := fsx.Exists(dst)
		if err != nil {
			return fmt.Errorf("failed to check for existence: %w", err)
		}

		// create the file if it doesn't exist
		if !exists {
			if _, err := fmt.Fprintf(progress, "[create]   %s\n", dst); err != nil {
				return fmt.Errorf("failed to report progress: %w", err)
			}
			if err := umaskfree.WriteFile(dst, []byte(content), umaskfree.DefaultFilePerm); err != nil {
				return fmt.Errorf("failed to write destination file: %w", err)
			}
		} else {
			if _, err := fmt.Fprintf(progress, "[skip]   %s\n", dst); err != nil {
				return fmt.Errorf("failed to report progress: %w", err)
			}
		}
	}

	// check that the stack can be loaded
	{
		if _, err := fmt.Fprintln(progress, "[checking]"); err != nil {
			return fmt.Errorf("failed to report progress: %w", err)
		}
		_, err := is.Project(ctx)
		if err != nil {
			return fmt.Errorf("failed to validate project: %w", err)
		}
	}

	return nil
}

const composeFileHeader = "# This file was automatically created and is updated by the distillery; DO NOT EDIT.\n\n"

// adds a header to the compose file.
func addComposeFileHeader(path string) (e error) {
	// read existing bytes
	bytes, err := os.ReadFile(path) // #nosec G304 -- intended
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// overwrite the file
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, umaskfree.DefaultFilePerm) // #nosec G304 -- intended
	if err != nil {
		return fmt.Errorf("failed to open compose file: %w", err)
	}
	defer errorsx.Close(f, &e, "file")

	// write the header
	if _, err := f.WriteString(composeFileHeader); err != nil {
		return nil
	}

	// write the original content
	if _, err := f.Write(bytes); err != nil {
		return fmt.Errorf("failed to write compose file; %w", err)
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
			bytes, err := os.ReadFile(path) // #nosec G304 -- intended
			if err != nil {
				return fmt.Errorf("unable to read existing file: %w", err)
			}

			// unmarshal it into a node, or bail out!
			node = new(yaml.Node)
			if err := yaml.Unmarshal(bytes, node); err != nil {
				return fmt.Errorf("unable to unmarshal existing file: %w", err)
			}
		case errors.Is(err, fs.ErrNotExist):
			// file does not exist => use default mode
			mode = umaskfree.DefaultFilePerm

			// use a nil existing node
			node = nil
		default:
			return fmt.Errorf("failed to stat file: %w", err)
		}
	}

	// update the node
	node, err := update(node)
	if err != nil {
		return fmt.Errorf("update function failed: %w", err)
	}

	// re-encode the bytes
	result, err := yaml.Marshal(node)
	if err != nil {
		return fmt.Errorf("failed to re-marshal: %w", err)
	}

	// write the bytes back!
	if err := umaskfree.WriteFile(path, result, mode); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

// writeEnvFile writes an environment file.
func writeEnvFile(path string, perm fs.FileMode, variables map[string]string) (e error) {
	// create the environment file
	file, err := umaskfree.Create(path, perm)
	if err != nil {
		return fmt.Errorf("failed to create env file: %w", err)
	}
	defer errorsx.Close(file, &e, "file")

	// write the file!
	_, err = dockerenv.WriteEnvFile(file, variables)
	if err != nil {
		return fmt.Errorf("failed to write env file: %w", err)
	}

	// and return nil
	return nil
}
