package validators

import (
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/pkg/errors"
)

func ValidateFile(env environment.Environment, path *string, dflt string) error {
	if *path == "" {
		*path = dflt
	}
	if !fsx.IsFile(env, *path) {
		return errors.Errorf("%q does not exist or is not a file", *path)
	}
	return nil
}

func ValidateDirectory(env environment.Environment, path *string, dflt string) error {
	if *path == "" {
		*path = dflt
	}
	if !fsx.IsDirectory(env, *path) {
		return errors.Errorf("%q does not exist or is not a directory", *path)
	}
	return nil
}
