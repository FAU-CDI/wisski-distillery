//spellchecker:words unpack
package unpack

//spellchecker:words path filepath github errors pkglib umaskfree
import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/pkg/errwrap"
	"github.com/tkw1536/pkglib/fsx/umaskfree"
)

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
		return fmt.Errorf("failed to install directory: %w", err)
	}

	// stat it!
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
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
func installResource(dst string, src string, fsys fs.FS, onInstallFile func(dst, src string)) (e error) {
	// open the srcFile
	srcFile, err := fsys.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open file to install: %w", err)
	}
	defer errwrap.Close(srcFile, "file to install", &e)

	// stat it!
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
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
			return fmt.Errorf("unable to create destination directory %q: %w", dst, err)
		}
	case dstErr != nil:
		return fmt.Errorf("unable to call stat on destination %q: %w", dst, dstErr)
	case !dstStat.IsDir():
		return fmt.Errorf("unable to open destination %q: %w", dst, errExpectedDirectoryButGotFile)
	}

	// NOTE(twiesing): We don't use fs.Walk here.
	// If we did, we'd have to reconstruct relative paths.
	// That would be very ugly!

	// read the directory
	entries, err := srcFile.ReadDir(-1)
	if err != nil {
		return fmt.Errorf("unable to read source directory %q: %w", srcFile, err)
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

func installFile(dst string, srcInfo fs.FileInfo, src fs.File) (e error) {
	// create the file using the right mode!
	file, err := umaskfree.Create(dst, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer errwrap.Close(file, "file", &e)

	// copy over the content!
	_, err = io.Copy(file, src)
	if err != nil {
		return fmt.Errorf("error writing to destination: %w", err)
	}
	return nil
}
