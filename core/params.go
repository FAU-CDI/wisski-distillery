package core

import (
	"os"
	"path/filepath"
)

// Params are used to initialize the excutable.
type Params struct {
	ConfigPath string // ConfigPath is the path to the configuration file for the distillery
}

// ParamsFromEnv creates a new set of parameters from the environment.
// Uses [ReadBaseDirectory] or falls back to [BaseDirectoryDefault].
func ParamsFromEnv() (params Params, err error) {

	// try to read the base directory!
	value, err := ReadBaseDirectory()
	switch {
	case os.IsNotExist(err):
		params.ConfigPath = BaseDirectoryDefault
	case err == nil:
		params.ConfigPath = value
	default:
		return params, err
	}

	// and add the configuration file name to it!
	params.ConfigPath = filepath.Join(params.ConfigPath, ConfigFile)
	return params, nil
}
