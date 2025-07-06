package cli

//spellchecker:words errors user path filepath strings github wisski distillery internal bootstrap pkglib umaskfree
import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"go.tkw01536.de/pkglib/fsx/umaskfree"
)

// metaConfigFile is the path to a configuration file that contains the path to the last used wdcli executable.
// It is expected to be in the current user's home directory.
//
// It should contain the path to a deployment directory.
const metaConfigFile = "." + bootstrap.Executable

// MetaConfigPath returns the full path to the MetaConfigPath().
func MetaConfigPath() (string, error) {
	// find the current user
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}
	return filepath.Join(usr.HomeDir, metaConfigFile), nil
}

var errReadBaseDirectoryEmpty = errors.New("`ReadBaseDirectory': directory is empty")

// ReadBaseDirectory reads the base deployment directory from the environment.
// Use [ParamsFromEnv] to initialize parameters completely.
//
// It does not perform any reading of files.
func ReadBaseDirectory() (value string, err error) {
	// get the path!
	path, err := MetaConfigPath()
	if err != nil {
		return "", err
	}

	// read the meta config file!
	contents, err := os.ReadFile(path) // #nosec G304 -- intended
	if err != nil {
		return "", fmt.Errorf("failed to read meta config file: %w", err)
	}

	// and trim the spaces!
	value = strings.TrimSpace(string(contents))

	// check that it is actually set!
	if len(value) == 0 {
		return "", errReadBaseDirectoryEmpty
	}

	// and return it!
	return value, nil
}

// WriteBaseDirectory writes the base directory to the environment, or returns an error.
func WriteBaseDirectory(dir string) error {
	// get the path!
	path, err := MetaConfigPath()
	if err != nil {
		return err
	}

	// just put the directory inside it!
	if err := umaskfree.WriteFile(path, []byte(dir), fs.ModePerm); err != nil {
		return fmt.Errorf("failed to create base directory: %w", err)
	}
	return nil
}
