// Package fsx provides convenient abstractions to work with the filesystem.
package fsx

import (
	"io/fs"

	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

// Exists checks if the given path exists
func Exists(env environment.Environment, path string) bool {
	_, err := env.Lstat(path)
	return err == nil
}

// IsDirectory checks if the provided path exists and is a directory
func IsDirectory(env environment.Environment, path string) bool {
	info, err := env.Stat(path)
	return err == nil && info.Mode().IsDir()
}

// IsFile checks if the provided path exists and is a regular file
func IsFile(env environment.Environment, path string) bool {
	info, err := env.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

// IsLink checks if the provided path exists and is a symlink
func IsLink(env environment.Environment, path string) bool {
	info, err := env.Lstat(path)
	return err == nil && info.Mode()&fs.ModeSymlink != 0
}
