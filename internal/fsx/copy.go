package fsx

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

var ErrCopySameFile = errors.New("src and dst must be different files")

// CopyFile copies a file from src to dst.
// When dst and src are the same file, returns ErrCopySameFile.
func CopyFile(dst, src string) error {
	if src == dst {
		return ErrCopySameFile
	}

	// open the source
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// stat it to get the mode!
	srcStat, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// open or create the destination
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, srcStat.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// and do the copy!
	_, err = io.Copy(dstFile, srcFile)
	return err
}

var ErrCopyNoDirectory = errors.New("dst is not a directory")

// CopyDirectory copies the directory src to dst recursively.
// The destination directory must exist, or an error is returned.
//
// onCopy, when not nil, is called for each file or directory being copied.
func CopyDirectory(dst, src string, onCopy func(dst, src string)) error {
	// TODO: Allow copying in parallel? Maybe with a mutex?

	// sanity checks
	if src == dst {
		return ErrCopySameFile
	}
	if !IsDirectory(dst) {
		return ErrCopyNoDirectory
	}

	// call onCopy for this directory!
	if onCopy != nil {
		onCopy(dst, src)
	}

	// iterate over the entries or bail out
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		eDest := filepath.Join(dst, name)
		eSrc := filepath.Join(src, name)

		// it is not a directory => Use CopyFile
		if !entry.IsDir() {
			if onCopy != nil {
				onCopy(eDest, eSrc)
			}

			// do the copy!
			if err := CopyFile(eDest, eSrc); err != nil {
				return err
			}

			continue
		}

		// find out the mode of the entry
		eInfo, err := entry.Info()
		if err != nil {
			return err
		}

		// make the target directory
		if err := os.Mkdir(eDest, eInfo.Mode()); err != nil {
			return err
		}

		// do the copy!
		if err := CopyDirectory(eDest, eSrc, onCopy); err != nil {
			return err
		}
	}

	return nil
}
