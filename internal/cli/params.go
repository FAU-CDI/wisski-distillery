package cli

//spellchecker:words context errors signal path filepath github wisski distillery internal bootstrap
import (
	"context"
	"errors"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
)

// Params are used to initialize the excutable.
//
//nolint:containedctx
type Params struct {
	ConfigPath string          // ConfigPath is the path to the configuration file for the distillery
	Context    context.Context // Context for the distillery
}

// ParamsFromEnv creates a new set of parameters from the environment.
// Uses [ReadBaseDirectory] or falls back to [BaseDirectoryDefault].
func ParamsFromEnv() (params Params, err error) {
	// try to read the base directory!
	value, err := ReadBaseDirectory()
	switch {
	case errors.Is(err, fs.ErrNotExist):
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
