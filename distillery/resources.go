// TODO: Rename this to resources oncen finished
package distillery

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// resourceEmbed contains all the resources required by the WissKI-Distillery package.
//go:embed all:resources
var resourceEmbed embed.FS

// InstallResource install a resource src into dest.
// When it encounters a directory, recursively installs the directory is called.
// For each installation item, onInstallFile is called, unless onInstallFile is nil.
//
// If src points to a file, dst must either be an existing file, or not exist.
// If src points to a directory, dst must either be an existing directory, or not exist.
func InstallResource(dst, src string, onInstallFile func(dst, src string)) error {
	return installFile(dst, resourceEmbed, src, onInstallFile)
}

var errExpectedFileButGotDirectory = errors.New("Expected a file, but got a directory")
var errExpectedDirectoryButGotFile = errors.New("Expected a directory, but got a file")

func installFile(dst string, fsys embed.FS, src string, onInstallFile func(dst, src string)) error {
	// call the on-install file path
	if onInstallFile != nil {
		onInstallFile(dst, src)
	}

	// open the source file!
	srcFile, err := fsys.Open(src)
	if err != nil {
		return errors.Wrapf(err, "Error opening source file %s", src)
	}
	defer srcFile.Close()

	// stat the source file to install
	srcStat, srcErr := srcFile.Stat()
	if srcErr != nil {
		return errors.Wrapf(srcErr, "Error calling stat on source %s", src)
	}

	// if it is a directory, we should recurse!
	if srcStat.IsDir() {
		return installDir(dst, srcStat, srcFile, fsys, src, onInstallFile)
	}

	// determine if we need to create the destination file, or if it already exists
	dstStat, dstErr := os.Stat(dst)
	flag := os.O_WRONLY
	switch {
	case os.IsNotExist(dstErr):
		flag |= os.O_CREATE
	case dstErr != nil:
		return errors.Wrapf(dstErr, "Error calling stat on destination %s", dst)
	case dstStat.IsDir():
		return errors.Wrapf(errExpectedFileButGotDirectory, "Error processing destination %s", dst)
	}

	// Open the file
	dstFile, err := os.OpenFile(dst, flag, srcStat.Mode())
	if err != nil {
		return errors.Wrapf(err, "Error opening destination %s", dst)
	}
	defer dstFile.Close()

	// copy over the content
	_, err = io.Copy(dstFile, srcFile)
	return errors.Wrapf(err, "Error writing to destination %s", dst)

}

func installDir(dst string, srcStat fs.FileInfo, srcFile fs.File, fsys embed.FS, src string, onInstallFile func(dst, src string)) error {
	// make sure it is a directory!
	dir, ok := srcFile.(fs.ReadDirFile)
	if !ok {
		return errExpectedDirectoryButGotFile
	}

	// create the destination
	dstStat, dstErr := os.Stat(dst)
	switch {
	case os.IsNotExist(dstErr):
		if err := os.MkdirAll(dst, srcStat.Mode()); err != nil {
			return errors.Wrapf(err, "Error creating destination directory %s", dst)
		}
	case dstErr != nil:
		return errors.Wrapf(dstErr, "Error calling stat on destination %s", dst)
	case !dstStat.IsDir():
		return errors.Wrapf(errExpectedDirectoryButGotFile, "Error opening destination %s", dst)
	case dstErr == nil:
	}

	// read the directory
	entries, err := dir.ReadDir(-1)
	if err != nil {
		return errors.Wrapf(err, "Error reading source directory %s", srcFile)
	}

	// iterate over all the children
	for _, entry := range entries {
		if err := func(dst, src string) error {
			return installFile(dst, fsys, src, onInstallFile)
		}(
			filepath.Join(dst, entry.Name()),
			filepath.Join(src, entry.Name()),
		); err != nil {
			return err
		}
	}

	return nil
}
