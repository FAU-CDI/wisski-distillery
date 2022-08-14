package fsx

import (
	"errors"
	"io"
	"os"
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
