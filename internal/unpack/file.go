package unpack

import (
	"io"
	"io/fs"
	"os"

	"github.com/pkg/errors"
)

// InstallFile installs the file from src into dst.
//
// If the destination path does not exist, it is created.
func InstallFile(dst string, src fs.File) error {
	// stat it!
	srcInfo, err := src.Stat()
	if err != nil {
		return err
	}

	// if this is a directory, something went wrong!
	if srcInfo.IsDir() {
		return errExpectedFileButGotDirectory
	}

	// and store it there!
	return installFile(dst, srcInfo, src)
}

func installFile(dst string, srcInfo fs.FileInfo, src fs.File) error {
	// create the file using the right mode!
	file, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer file.Close()

	// copy over the content!
	_, err = io.Copy(file, src)
	return errors.Wrapf(err, "Error writing to destination %s", dst)
}

// InstallTemplate unpacks the resource located at src in fsys, then processes it as a template, and eventually writes it to dst.
// Any existing file is truncated and overwritten.
//
// See [WriteTemplate] for possible errors.
func InstallTemplate(dst string, context map[string]string, src string, fsys fs.FS) error {

	// open the srcFile
	srcFile, err := fsys.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// stat it
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// check if it is a directory
	if srcInfo.IsDir() {
		return errExpectedFileButGotDirectory
	}

	// open the destination file
	file, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer file.Close()

	// write the file!
	return WriteTemplate(file, context, srcFile)
}
