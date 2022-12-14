package environment

import (
	"context"
	"io"
	"io/fs"
	"net"
	"time"

	"github.com/tkw1536/goprogram/stream"
)

// Environment represents an environment that a program can run it.
// It mostly mimics the interfaces of the [os] package.
type Environment interface {
	isEnv()

	GetEnv(name string) string

	Stat(path string) (fs.FileInfo, error)
	Lstat(path string) (fs.FileInfo, error)

	Readlink(path string) (string, error)
	Symlink(oldname, newname string) error

	ReadDir(name string) ([]fs.DirEntry, error)

	Open(path string) (fs.File, error)
	Chtimes(name string, atime time.Time, mtime time.Time) error
	SameFile(f1, f2 fs.FileInfo) bool

	Create(path string, mode fs.FileMode) (WritableFile, error)

	Mkdir(path string, mode fs.FileMode) error
	MkdirAll(path string, mode fs.FileMode) error

	Remove(path string) error
	RemoveAll(path string) error

	WalkDir(root string, fn fs.WalkDirFunc) error

	Abs(path string) (string, error)

	Listen(network, address string) (net.Listener, error)
	DialContext(context context.Context, network, address string) (net.Conn, error)

	Executable() (string, error)
	Exec(ctx context.Context, io stream.IOStream, workdir string, exe string, argv ...string) func() int
	LookPathAbs(name string) (string, error)
}

type WritableFile interface {
	fs.File
	io.Writer
}

func init() {
	var _ Environment = new(Native)
}
