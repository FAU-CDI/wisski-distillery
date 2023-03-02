// Package fsx provides additional file system functionality.
//
// Several functions in this package provide umask-ignoring functions.
// Using these functions intervenes with the global umask.
//
// It is not safe to use functions provided by the standard go library concurrently with this function.
// Users should take care that no other code in their application uses these functions.
package fsx

import "io/fs"

// DefaultFilePerm should be used by callers to use a consistent file mode for new files.
const DefaultFilePerm fs.FileMode = 0666

// DefaultDirPerm should be used by callers to use a consistent mode for new directories.
const DefaultDirPerm fs.FileMode = fs.ModeDir | fs.ModePerm
