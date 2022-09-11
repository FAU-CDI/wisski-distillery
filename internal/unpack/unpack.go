// Package unpack unpacks files and templates to a target directory
package unpack

import (
	"bytes"
	"errors"
	"io/fs"
)

var errExpectedFileButGotDirectory = errors.New("expected a file, but got a directory")
var errExpectedDirectoryButGotFile = errors.New("expected a directory, but got a file")

// InstallResource installs the resource at src within fsys to dst.
//
// OnInstallFile is called for each source and destination file.
// OnInstallFile may be nil.
//
// See [InstallDir] or [InstallFile].
func InstallResource(dst string, src string, fsys fs.FS, onInstallFile func(dst, src string)) error {
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

// UnpackTemplate unpacks the given file template and template.
// See [WriteTemplate] for possible errors.
func UnpackTemplate(context map[string]string, src fs.File) ([]byte, fs.FileMode, error) {
	// stat the source file to install
	srcStat, srcErr := src.Stat()
	if srcErr != nil {
		return nil, 0, srcErr
	}

	// should not be a directory
	if srcStat.IsDir() {
		return nil, 0, errExpectedFileButGotDirectory
	}

	// read all the bytes into a buffer
	var buffer bytes.Buffer
	err := WriteTemplate(&buffer, context, src)
	return buffer.Bytes(), srcStat.Mode(), err
}
