// Package targz provides facilities for packaging tar.gz files
//
//spellchecker:words targz
package targz

//spellchecker:words archive compress gzip path filepath pkglib errorsx umaskfree
import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/fsx/umaskfree"
)

// Package packages the source directory into a 'tar.gz' file into destination.
// If the destination already exists, it is truncated.
//
// onCopy, when not nil, is called for each file being copied into the archive.
func Package(dst, src string, onCopy func(rel string, src string)) (count int64, e error) {
	// create the target archive
	archive, err := umaskfree.Create(dst, umaskfree.DefaultFilePerm)
	if err != nil {
		return 0, fmt.Errorf("failed to create file: %w", err)
	}
	defer errorsx.Close(archive, &e, "archive file")

	// create a gzip writer
	zipHandle := gzip.NewWriter(archive)
	defer errorsx.Close(zipHandle, &e, "zip handle")

	// create a tar writer
	tarHandle := tar.NewWriter(zipHandle)
	defer errorsx.Close(tarHandle, &e, "tar handle")

	// and walk through it!
	e = filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) (e error) {
		if err != nil {
			return err
		}

		// determine the relative path
		var relpath string
		relpath, err = filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %q: %w", path, err)
		}

		// call the oncopy!
		if onCopy != nil {
			onCopy(relpath, path)
		}

		// read mode etc
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("failed to get file info for %q: %w", path, err)
		}

		// FIXME: How do we handle

		// create a file info header!
		tInfo, err := tar.FileInfoHeader(info, relpath)
		if err != nil {
			return fmt.Errorf("failed to create info header for %q: %w", path, err)
		}
		tInfo.Name = filepath.ToSlash(relpath)

		// write it!
		if err := tarHandle.WriteHeader(tInfo); err != nil {
			return fmt.Errorf("failed to write tar header for %q: %w", path, err)
		}

		// if it's not a regular file, we are done
		if !entry.Type().IsRegular() {
			return nil
		}

		// open the file
		handle, err := os.Open(path) // #nosec G304 -- intended
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer errorsx.Close(handle, &e, "file")

		// and copy it into the archive
		ccount, err := io.Copy(tarHandle, handle)
		count += ccount
		if err != nil {
			return fmt.Errorf("failed to copy %q into archive: %w", path, err)
		}
		return nil
	})
	return
}
