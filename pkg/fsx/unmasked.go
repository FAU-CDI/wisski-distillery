package fsx

import (
	"io/fs"
	"os"
	"sync"
	"syscall"
	"time"
)

// mask is the global mask lock
var m mask

// mask allows disabling and re-enabling the global umask.
// it is used by allow functions of this package.
type mask struct {
	l     sync.Mutex // locked?
	umask int        // previous mask
}

// Lock blocks until no other function is using this umask
// and then sets it to 0.
func (mask *mask) Lock() {
	mask.l.Lock()
	mask.umask = syscall.Umask(0)
}

func (mask *mask) Unlock() {
	mask.umask = syscall.Umask(mask.umask)
	mask.l.Unlock()
}

// WriteFile is like [os.WriteFile], but ignores the umask.
func WriteFile(path string, data []byte, mode fs.FileMode) error {
	handle, err := Create(path, mode)
	if err != nil {
		return err
	}
	defer handle.Close()

	if _, err := handle.Write(data); err != nil {
		return err
	}

	return nil
}

// Create creates a new file with the given mode.
// This function ignores the umask.
func Create(path string, mode fs.FileMode) (*os.File, error) {
	m.Lock()
	defer m.Unlock()

	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
}

// Mkdir creates a new directory with the given mode.
// This function ignores the umask.
func Mkdir(path string, mode fs.FileMode) error {
	m.Lock()
	defer m.Unlock()

	return os.Mkdir(path, fs.ModeDir|mode)
}

// MkdirAll creates a new directory and all potentially missing parent directories.
// This function ignores the umask.
func MkdirAll(path string, mode fs.FileMode) error {
	m.Lock()
	defer m.Unlock()

	return os.MkdirAll(path, fs.ModeDir|mode)
}

// Touch touches a file.
// It is similar to the unix 'touch' command.
//
// If the file does not exist exists, it is created using [Create].
// If the file does exist, it's access and modification times are updated to the current time.
//
// This function ignores the umask.
func Touch(path string, perm fs.FileMode) error {
	if perm == 0 {
		perm = DefaultFilePerm
	}
	_, err := os.Stat(path)
	switch {
	case os.IsNotExist(err):
		f, err := Create(path, perm)
		if err != nil {
			return err
		}
		defer f.Close()
		return nil
	case err != nil:
		return err
	default:
		now := time.Now().Local()
		return os.Chtimes(path, now, now)
	}
}
