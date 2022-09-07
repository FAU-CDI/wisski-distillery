package targz

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Package packages the directory src into dst.
// onCopy, when not nil, is called for each file being copied into the archive.
func Package(dst, src string, onCopy func(rel string, src string)) (count int64, err error) {
	// create the target archive
	archive, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer archive.Close()

	// create a gzip writer
	zipHandle := gzip.NewWriter(archive)
	defer zipHandle.Close()

	// create a tar writer
	tarHandle := tar.NewWriter(zipHandle)
	defer tarHandle.Close()

	// and walk through it!
	err = filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// determine the relative path
		var relpath string
		relpath, err = filepath.Rel(src, path)
		if err != nil {
			return err
		}

		if onCopy != nil {
			onCopy(relpath, path)
		}

		// create a file info header!
		tInfo, err := tar.FileInfoHeader(info, relpath)
		if err != nil {
			return err
		}
		tInfo.Name = filepath.ToSlash(relpath)

		// write it!
		if err := tarHandle.WriteHeader(tInfo); err != nil {
			return err
		}

		// a directory => no more writing required
		if info.IsDir() {
			return nil
		}

		// open the file
		handle, err := os.Open(path)
		if err != nil {
			return err
		}
		defer handle.Close()

		// and copy it into the archive
		ccount, err := io.Copy(tarHandle, handle)
		count += ccount
		return err
	})
	return
}
