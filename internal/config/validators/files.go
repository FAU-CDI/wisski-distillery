package validators

import (
	"github.com/pkg/errors"
	"github.com/tkw1536/pkglib/fsx"
)

func ValidateFile(path *string, dflt string) error {
	if *path == "" {
		*path = dflt
	}
	if !fsx.IsRegular(*path) {
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
