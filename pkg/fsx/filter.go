package fsx

import (
	"io/fs"
	"path/filepath"

	"github.com/tkw1536/goprogram/lib/collection"
)

// Censor returns a new filesystem censors files for which the censor function returns true.
//
// A censored file cannot be opened by the filesystem and return [fs.ErrNotExist].
// Hard and Soft Links pointing to the file might still read it.
func Censor(fsys fs.FS, censor func(name string) bool) fs.FS {
	return &censorFS{
		fsys:   fsys,
		censor: censor,
	}
}

type censorFS struct {
	censor func(path string) bool
	fsys   fs.FS
}

func (cf *censorFS) Sub(path string) (fs.FS, error) {
	sub, err := fs.Sub(cf.fsys, path)
	if err != nil {
		return nil, err
	}
	return &censorFS{
		censor: func(name string) bool {
			return cf.censor(filepath.Join(path, name))
		},
		fsys: sub,
	}, nil
}

func (ef *censorFS) Open(name string) (fs.File, error) {
	if ef.censor(name) {
		return nil, fs.ErrNotExist
	}

	file, err := ef.fsys.Open(name)

	// we need to also censor the ReadDir function of the returned file
	// this is to prevent the file from appearing in directory listings.
	if rdf, ok := file.(fs.ReadDirFile); ok {
		return &censorFSFile{
			ReadDirFile: rdf,

			name:   name,
			censor: ef,
		}, err
	}
	return file, err
}

type censorFSFile struct {
	fs.ReadDirFile

	name   string
	censor *censorFS
}

func (f *censorFSFile) ReadDir(n int) ([]fs.DirEntry, error) {
	entries, err := f.ReadDirFile.ReadDir(n)
	return f.censor.handleReadDir(f.name, entries, err)
}

func (ef *censorFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if ef.censor(name) {
		return nil, fs.ErrNotExist
	}

	// censor ReadDir() entries too
	entries, err := fs.ReadDir(ef.fsys, name)
	return ef.handleReadDir(name, entries, err)
}

// handleReadDir censors a ReadDir call
func (ef *censorFS) handleReadDir(base string, entries []fs.DirEntry, err error) ([]fs.DirEntry, error) {
	entries = collection.Filter(entries, func(entry fs.DirEntry) bool {
		return !ef.censor(filepath.Join(base, entry.Name()))
	})
	return entries, err
}

func (ef *censorFS) ReadFile(name string) ([]byte, error) {
	if ef.censor(name) {
		return nil, fs.ErrNotExist
	}
	return fs.ReadFile(ef.fsys, name)
}
