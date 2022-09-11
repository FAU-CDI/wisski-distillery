package fsx

import "io/fs"

// OpenFS opens the named file in filesystem.
// If opening the file results in an error, returns [ErrFile].
func OpenFS(name string, fsys fs.FS) fs.File {
	file, err := fsys.Open(name)
	if err != nil {
		return ErrFile{Err: err}
	}
	return file
}

// ErrFile implements a no-op [fs.File].
//
// Every operation will return an underlying error
type ErrFile struct {
	Err error
}

func (err ErrFile) Stat() (fs.FileInfo, error) {
	return nil, err.Err
}
func (err ErrFile) Read([]byte) (int, error) {
	return 0, err.Err
}

func (err ErrFile) Close() error {
	return err.Err
}
