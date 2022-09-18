package component

import (
	"io"

	"github.com/tkw1536/goprogram/stream"
)

// Backupable represents a component with a Backup method
type Backupable interface {
	Component

	// BackupName returns a new name to be used as an argument for path.
	BackupName() string

	// Backup backs up this component into the destination path path
	Backup(context BackupContext) error
}

// BackupContext is the context for backups
type BackupContext interface {
	// IO returns the input output stream belonging to this backup file
	IO() stream.IOStream

	// Name creates a new directory inside the destination.
	// Passing the empty path creates the destination as a directory.
	//
	// It then allows op to fill the file.
	AddDirectory(path string, op func() error) error

	// CopyFile copies a file from src to dst.
	CopyFile(dst, src string) error

	// AddFile creates a new file at the provided path inside the destination.
	// Passing the empty path creates the destination as a file.
	//
	// It then allows op to write to the file.
	//
	// The op function must not retain file.
	// The underlying file does not need to be closed.
	// AddFile will not return before op has returned.
	AddFile(path string, op func(file io.Writer) error) error
}
