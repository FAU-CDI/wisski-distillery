//spellchecker:words validators
package validators

//spellchecker:words github errors pkglib
import (
	"github.com/pkg/errors"
	"github.com/tkw1536/pkglib/fsx"
)

func ValidateFile(path *string, dflt string) error {
	if *path == "" {
		*path = dflt
	}
	isFile, err := fsx.IsRegular(*path, true)
	if err != nil {
		return err
	}
	if !isFile {
		return errors.Errorf("%q does not exist or is not a file", *path)
	}
	return nil
}

func ValidateDirectory(path *string, dflt string) error {
	if *path == "" {
		*path = dflt
	}
	isDirectory, err := fsx.IsDirectory(*path, true)
	if err != nil {
		return err
	}
	if !isDirectory {
		return errors.Errorf("%q does not exist or is not a directory", *path)
	}
	return nil
}
