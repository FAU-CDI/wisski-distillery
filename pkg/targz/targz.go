// Package targz provides facilities for packaging tar.gz files
//
//spellchecker:words targz
package targz

//spellchecker:words archive compress gzip path filepath github pkglib umaskfree
import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tkw1536/pkglib/fsx/umaskfree"
)

// Package packages the source directory into a 'tar.gz' file into destination.
// If the destination already exists, it is truncated.
//
// onCopy, when not nil, is called for each file being copied into the archive.
func Package(dst, src string, onCopy func(rel string, src string)) (count int64, err error) {
	// create the target archive
	archive, err := umaskfree.Create(dst, umaskfree.DefaultFilePerm)
	if err != nil {
		return 0, fmt.Errorf("failed to create file: %w", err)
	}
	defer archive.Close()

	// create a gzip writer
	zipHandle := gzip.NewWriter(archive)
	defer zipHandle.Close()

	// create a tar writer
	tarHandle := tar.NewWriter(zipHandle)
	defer tarHandle.Close()

	// and walk through it!
	err = filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) error {
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
		defer handle.Close()

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
