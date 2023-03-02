package fsx

import (
	"io/fs"
	"os"
)

// DefaultDirPerm should be used by callers to use a consistent mode for new directories.
const DefaultDirPerm fs.FileMode = fs.ModeDir | fs.ModePerm

// Mkdir is like [os.Mkdir].
func Mkdir(path string, mode fs.FileMode) error {
	m.Lock()
	defer m.Unlock()

	return os.Mkdir(path, fs.ModeDir|mode)
}

// MkdirAll is like [os.MkdirAll].
func MkdirAll(path string, mode fs.FileMode) error {
	m.Lock()
	defer m.Unlock()

	return os.MkdirAll(path, fs.ModeDir|mode)
}
