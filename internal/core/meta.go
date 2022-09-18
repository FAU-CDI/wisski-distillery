package core

import (
	"errors"
	"io/fs"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

// MetaConfigFile is the path to a configuration file that contains the path to the last used wdcli executable.
// It is expected to be in the current user's home directory.
//
// You probably want to use [MetaConfigPath] instead.
//
// It should contain the path to a deployment directory.
const MetaConfigFile = "." + Executable

// MetaConfigPath returns the full path to the MetaConfigPath()
func MetaConfigPath() (string, error) {
	// find the current user
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, MetaConfigFile), nil
}

var errReadBaseDirectoryEmpty = errors.New("ReadBaseDirectory: Directory is empty")

// ReadBaseDirectory reads the base deployment directory from the environment.
// Use [ParamsFromEnv] to initialize parameters completely.
//
// It does not perform any reading of files.
func ReadBaseDirectory(env environment.Environment) (value string, err error) {
	// get the path!
	path, err := MetaConfigPath()
	if err != nil {
		return "", err
	}

	// read the meta config file!
	contents, err := environment.ReadFile(env, path)
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
func WriteBaseDirectory(env environment.Environment, dir string) error {
	// get the path!
	path, err := MetaConfigPath()
	if err != nil {
		return err
	}

	// just put the directory inside it!
	return environment.WriteFile(env, path, []byte(dir), fs.ModePerm)
}
