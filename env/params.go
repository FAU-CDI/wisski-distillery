package env

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/tkw1536/goprogram/exit"
)

// Params are parameters used for initialization of the environment
type Params struct {
	BaseDirectory string
}

// ConfigFilePath returns the path to the configuration file
func (params Params) ConfigFilePath() string {
	if params.BaseDirectory == "" {
		return ""
	}
	return filepath.Join(params.BaseDirectory, ".env")
}

var errUnableToLoadParams = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "Unable to configure wdcli environment: %s",
}

const BaseDirectoryDefault = "/var/www/deploy"

// ParamsFromEnv creates a new set of parameters from the environment.
// There is no guarantee that the parameters are valid.
func ParamsFromEnv() (params Params, err error) {
	// try to read the base directory
	value, err := ReadBaseDirectory()
	switch {
	case os.IsNotExist(err):
		params.BaseDirectory = BaseDirectoryDefault
	case err == nil:
		params.BaseDirectory = value
	default:
		return params, errUnableToLoadParams.WithMessageF(err)
	}

	return params, nil
}

var baseConfigFile = ".wdcli"

// ReadBaseDirectory reads the base directory from the environment, or an empty string
func ReadBaseDirectory() (value string, err error) {
	// find the current user
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	// read the base config file!
	contents, err := os.ReadFile(filepath.Join(usr.HomeDir, baseConfigFile))
	if err != nil {
		return "", err
	}

	// and trim the spaces!
	value = strings.TrimSpace(string(contents))

	// check that it is actually set!
	if len(value) == 0 {
		return "", errors.New("ReadBaseDirectory: Directory is empty")
	}

	// and return it!
	return value, nil
}

// WriteBaseDirectory writes the base directory to the environment, or returns an error
func WriteBaseDirectory(dir string) error {
	// find the current user
	usr, err := user.Current()
	if err != nil {
		return err
	}

	// read the base config file!
	return os.WriteFile(
		filepath.Join(usr.HomeDir, baseConfigFile),
		[]byte(dir),
		os.ModePerm,
	)
}
