package environment

import (
	"bytes"
	"io"
	"io/fs"
	"os"

	"github.com/tkw1536/goprogram/stream"
)

// ExecCommandError is returned by Exec when a command could not be executed.
// This typically hints that the executable cannot be found, but may have other causes.
const ExecCommandError = 127

// DefaultFilePerm is the default mode to use for files
const DefaultFilePerm fs.FileMode = 0666

// DefaultDirPerm is the default mode to use for directories
const DefaultDirPerm fs.FileMode = os.ModeDir & fs.ModePerm

// IsExist checks if the provided error represents a 'does not exist' errror
func IsExist(err error) bool {
	return os.IsExist(err)
}

// IsNotExist checks if the provided error represents a 'does exist' error
func IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

// WriteFile is like [os.WriteFile].
func WriteFile(env Environment, path string, data []byte, mode fs.FileMode) error {
	handle, err := env.Create(path, mode)
	if err != nil {
		return err
	}
	defer handle.Close()

	if _, err := handle.Write(data); err != nil {
		return err
	}

	return nil
}

// ReadFile is like [os.ReadFile]
func ReadFile(env Environment, path string) ([]byte, error) {
	// open the file!
	file, err := env.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// copy everything into a buffer!
	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, file); err != nil {
		return nil, err
	}

	// return the buffer contents!
	return buffer.Bytes(), nil
}

// MustExec is like Exec, except that it returns true if the command exited successfully, and else false.
func MustExec(env Environment, io stream.IOStream, workdir string, exe string, argv ...string) bool {
	return env.Exec(io, workdir, exe, argv...) == 0
}
