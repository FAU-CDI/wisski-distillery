package cli

import (
	"errors"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
)

// metaConfigFile is the path to a configuration file that contains the path to the last used wdcli executable.
// It is expected to be in the current user's home directory.
//
// It should contain the path to a deployment directory.
const metaConfigFile = "." + bootstrap.Executable

// MetaConfigPath returns the full path to the MetaConfigPath()
func MetaConfigPath() (string, error) {
	// find the current user
	usr, err := user.Current()
	if err != nil {
		return "", err
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
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
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

// WriteBaseDirectory writes the base directory to the environment, or returns an error
func WriteBaseDirectory(dir string) error {
	// get the path!
	path, err := MetaConfigPath()
	if err != nil {
		return err
	}

	// just put the directory inside it!
	return fsx.WriteFile(path, []byte(dir), fs.ModePerm)
}
