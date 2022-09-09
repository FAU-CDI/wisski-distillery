package embed

import (
	"os"

	"github.com/FAU-CDI/wisski-distillery/internal/unpack"
	"github.com/pkg/errors"
)

// InstallTemplates open the resource src, and installs it into dst.
// the template resource must fit into memory.
//
// For each variable ${THING} inside dest, a key 'THING' must exist in context.
// Extra or missing template keys are an error.
func InstallTemplate(dst, src string, context map[string]string) error {
	// open the source file!
	srcFile, err := resourceEmbed.Open(src)
	if err != nil {
		return errors.Wrapf(err, "Error opening source file %s", src)
	}
	defer srcFile.Close()

	// write the template
	bytes, srcMode, err := unpack.UnpackTemplate(context, srcFile)
	if err != nil {
		return err
	}

	// determine if we need to create the destination file, or if it already exists
	dstStat, dstErr := os.Stat(dst)
	switch {
	case os.IsNotExist(dstErr):
	case dstErr != nil:
		return errors.Wrapf(dstErr, "Error calling stat on destination %s", dst)
	case dstStat.IsDir():
		return errors.Wrapf(errExpectedFileButGotDirectory, "Error processing destination %s", dst)
	}

	// open and write the destination file
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcMode)
	if err != nil {
		return errors.Wrapf(err, "Unable to open file %s", dst)
	}
	_, err = dstFile.Write(bytes)
	return errors.Wrapf(err, "Unable to write destination %s", dst)
}

// ReadTemplate is like InstallTemplate, except that it writes template into a byte slice and returns it.
func ReadTemplate(src string, context map[string]string) ([]byte, error) {
	// open the source file!
	srcFile, err := resourceEmbed.Open(src)
	if err != nil {
		return nil, errors.Wrapf(err, "Error opening source file %s", src)
	}
	defer srcFile.Close()

	// and return it
	bytes, _, err := unpack.UnpackTemplate(context, srcFile)
	return bytes, err
}
