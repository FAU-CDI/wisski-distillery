// Package fsx provides convenient abstractions to work with the filesystem.
package fsx

import (
	"os"
)

// Exists checks if the given path exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDirectory checks if the provided path exists and is a directory
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsDir()
}

// IsFile checks if the provided path exists and is a regular file
func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}
