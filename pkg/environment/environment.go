package environment

import (
	"io"
	"io/fs"
)

// Environment represents an environment that a program can run it.
// It mostly mimics the interfaces of the [os] package.
type Environment interface {
	isEnv()

	Create(path string, mode fs.FileMode) (WritableFile, error)
	Mkdir(path string, mode fs.FileMode) error
	MkdirAll(path string, mode fs.FileMode) error
}

type WritableFile interface {
	fs.File
	io.Writer
}

func init() {
	var _ Environment = new(Native)
}
