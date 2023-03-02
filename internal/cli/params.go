package cli

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

// Params are used to initialize the excutable.
type Params struct {
	ConfigPath string          // ConfigPath is the path to the configuration file for the distillery
	Context    context.Context // Context for the distillery
}

// ParamsFromEnv creates a new set of parameters from the environment.
// Uses [ReadBaseDirectory] or falls back to [BaseDirectoryDefault].
func ParamsFromEnv() (params Params, err error) {
	var native environment.Environment

	// try to read the base directory!
	value, err := ReadBaseDirectory(native) // TODO: Are we sure about the native environment here?
	switch {
	case os.IsNotExist(err):
		params.ConfigPath = bootstrap.BaseDirectoryDefault
	case err == nil:
		params.ConfigPath = value
	default:
		return params, err
	}

	// and add the configuration file name to it!
	params.ConfigPath = filepath.Join(params.ConfigPath, bootstrap.ConfigFile)

	// generate a new context
	params.Context, _ = signal.NotifyContext(context.Background(), os.Interrupt)

	// and return the params!
	return params, nil
}
