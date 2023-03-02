package validators

import (
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/pkg/errors"
)

func ValidateFile(path *string, dflt string) error {
	if *path == "" {
		*path = dflt
	}
	if !fsx.IsFile(*path) {
		return errors.Errorf("%q does not exist or is not a file", *path)
	}
	return nil
}

func ValidateDirectory(path *string, dflt string) error {
	if *path == "" {
		*path = dflt
	}
	if !fsx.IsDirectory(*path) {
		return errors.Errorf("%q does not exist or is not a directory", *path)
	}
	return nil
}
