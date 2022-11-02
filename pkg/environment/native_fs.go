package environment

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

type Native struct {
	ulock sync.Mutex
	umask int
}

func (*Native) isEnv() {}

func (n *Native) setMask(umask int) {
	n.ulock.Lock()
	n.umask = syscall.Umask(umask)
}

func (n *Native) resetMask() {
	syscall.Umask(n.umask)
	n.ulock.Unlock()
}

func (*Native) GetEnv(name string) string {
	return os.Getenv(name)
}

func (*Native) Stat(path string) (fs.FileInfo, error) {
	return os.Stat(path)
}

func (*Native) Lstat(path string) (fs.FileInfo, error) {
	return os.Lstat(path)
}

func (*Native) Readlink(path string) (string, error) {
	return os.Readlink(path)
}

func (*Native) Symlink(oldname, newname string) error {
	return os.Symlink(oldname, newname)
}

func (*Native) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

func (*Native) SameFile(f1, f2 fs.FileInfo) bool {
	return os.SameFile(f1, f2)
}

func (*Native) WalkDir(path string, f fs.WalkDirFunc) error {
	return filepath.WalkDir(path, f)
}

func (*Native) Executable() (string, error) {
	return os.Executable() // TODO: not sure this works with the remote concepts
}

func (*Native) Open(path string) (fs.File, error) {
	return os.Open(path)
}

func (n *Native) Create(path string, mode fs.FileMode) (WritableFile, error) {
	n.ulock.Lock()
	defer n.ulock.Unlock()

	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
}

func (*Native) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(name, atime, mtime)
}

func (n *Native) Mkdir(path string, mode fs.FileMode) error {
	n.setMask(0)
	defer n.resetMask()

	return os.Mkdir(path, fs.ModeDir|mode)
}

func (n *Native) MkdirAll(path string, mode fs.FileMode) error {
	n.setMask(0)
	defer n.resetMask()

	return os.MkdirAll(path, fs.ModeDir|mode)
}

func (*Native) Remove(path string) error {
	return os.Remove(path)
}

func (*Native) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (*Native) Abs(path string) (string, error) {
	return filepath.Abs(path)
}
