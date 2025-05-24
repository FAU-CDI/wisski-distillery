//spellchecker:words validators
package validators

//spellchecker:words github pkglib
import (
	"fmt"
	"io/fs"

	"github.com/tkw1536/pkglib/fsx"
)

func ValidateFile(path *string, dflt string) error {
	if *path == "" {
		*path = dflt
	}
	isFile, err := fsx.IsRegular(*path, true)
	if err != nil {
		return fmt.Errorf("failed to check for regular file: %w", err)
	}
	if !isFile {
		return fmt.Errorf("%q does not exist or is not a file: %w", *path, fs.ErrNotExist)
	}
	return nil
}

func ValidateDirectory(path *string, dflt string) error {
	if *path == "" {
		*path = dflt
	}
	isDirectory, err := fsx.IsDirectory(*path, true)
	if err != nil {
		return fmt.Errorf("failed to check for directory: %w", err)
	}
	if !isDirectory {
		return fmt.Errorf("%q does not exist or is not a directory: %w", *path, fs.ErrNotExist)
	}
	return nil
}
