package unpack

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// InstallDir installs the directory at src within fsys to dst.
//
// onInstallFile is called for each file or directory being installed.
//
// If the destination path does not exist, it is created using [os.MakeDirs]
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

func installDir(dst string, srcInfo fs.FileInfo, srcFile fs.ReadDirFile, src string, fsys fs.FS, onInstallFile func(dst, src string)) error {
	// create the destination
	dstStat, dstErr := os.Stat(dst)
	switch {
	case os.IsNotExist(dstErr):
		if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
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
		if err := InstallResource(
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
