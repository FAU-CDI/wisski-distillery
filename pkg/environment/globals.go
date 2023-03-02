package environment

import (
	"io"
	"io/fs"
	"os"

	"github.com/tkw1536/pkglib/pools"
)

// DefaultFilePerm is the default mode to use for files
const DefaultFilePerm fs.FileMode = 0666

// DefaultDirPerm is the default mode to use for directories
const DefaultDirPerm fs.FileMode = fs.ModeDir | fs.ModePerm

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
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// copy everything into a buffer!
	buffer := pools.GetBuffer()
	defer pools.ReleaseBuffer(buffer)

	if _, err := io.Copy(buffer, file); err != nil {
		return nil, err
	}

	// return the buffer contents!
	return buffer.Bytes(), nil
}
