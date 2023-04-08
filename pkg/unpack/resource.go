package unpack

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/tkw1536/pkglib/fsx/umaskfree"
)

var errExpectedFileButGotDirectory = errors.New("expected a file, but got a directory")
var errExpectedDirectoryButGotFile = errors.New("expected a directory, but got a file")

// InstallDir installs the directory at src within fsys to dst.
//
// onInstallFile is called for each file or directory being installed.
//
// If the destination path does not exist, it is created using [environment.MakeDirs]
// The directory is installed recursively.
func InstallDir(dst string, src string, fsys fs.FS, onInstallFile func(dst, src string)) error {
	// open the source file
	srcFile, err := fsys.Open(src)
	if err != nil {
		return err
	}

	// stat it!
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// make sure it's a file!
	if !srcInfo.IsDir() {
		return errExpectedDirectoryButGotFile
	}

	// call the hook (if any)
	if onInstallFile != nil {
		onInstallFile(dst, src)
	}

	// do the installation of the directory.
	// the type cast should be safe.
	return installDir(dst, srcInfo, srcFile.(fs.ReadDirFile), src, fsys, onInstallFile)
}

// installResource installs the resource at src within fsys to dst.
//
// OnInstallFile is called for each source and destination file.
// OnInstallFile may be nil.
func installResource(dst string, src string, fsys fs.FS, onInstallFile func(dst, src string)) error {
	// open the srcFile
	srcFile, err := fsys.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// stat it!
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// call the hook (if any)
	if onInstallFile != nil {
		onInstallFile(dst, src)
	}

	// this is a directory, so the cast is safe!
	if srcInfo.IsDir() {
		return installDir(dst, srcInfo, srcFile.(fs.ReadDirFile), src, fsys, onInstallFile)
	}

	// this is a regular file!
	return installFile(dst, srcInfo, srcFile)
}

func installDir(dst string, srcInfo fs.FileInfo, srcFile fs.ReadDirFile, src string, fsys fs.FS, onInstallFile func(dst, src string)) error {
	// create the destination
	dstStat, dstErr := os.Stat(dst)
	switch {
	case errors.Is(dstErr, fs.ErrNotExist):
		if err := umaskfree.MkdirAll(dst, srcInfo.Mode()); err != nil {
			return errors.Wrapf(err, "Error creating destination directory %s", dst)
		}
	case dstErr != nil:
		return errors.Wrapf(dstErr, "Error calling stat on destination %s", dst)
	case !dstStat.IsDir():
		return errors.Wrapf(errExpectedDirectoryButGotFile, "Error opening destination %s", dst)
	}

	// NOTE(twiesing): We don't use fs.Walk here.
	// If we did, we'd have to reconstruct relative paths.
	// That would be very ugly!

	// read the directory
	entries, err := srcFile.ReadDir(-1)
	if err != nil {
		return errors.Wrapf(err, "Error reading source directory %s", srcFile)
	}

	// iterate over all the children
	for _, entry := range entries {
		if err := installResource(
			filepath.Join(dst, entry.Name()),
			filepath.Join(src, entry.Name()),
			fsys,
			onInstallFile,
		); err != nil {
			return err
		}
	}

	return nil
}

func installFile(dst string, srcInfo fs.FileInfo, src fs.File) error {
	// create the file using the right mode!
	file, err := umaskfree.Create(dst, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer file.Close()

	// copy over the content!
	_, err = io.Copy(file, src)
	return errors.Wrapf(err, "Error writing to destination %s", dst)
}
