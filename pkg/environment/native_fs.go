package environment

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type Native struct{}

func (Native) isEnv() {}

func (Native) GetEnv(name string) string {
	return os.Getenv(name)
}

func (Native) Stat(path string) (fs.FileInfo, error) {
	return os.Stat(path)
}

func (Native) Lstat(path string) (fs.FileInfo, error) {
	return os.Lstat(path)
}

func (Native) Readlink(path string) (string, error) {
	return os.Readlink(path)
}

func (Native) Symlink(oldname, newname string) error {
	return os.Symlink(oldname, newname)
}

func (Native) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

func (Native) SameFile(f1, f2 fs.FileInfo) bool {
	return os.SameFile(f1, f2)
}

func (Native) WalkDir(path string, f fs.WalkDirFunc) error {
	return filepath.WalkDir(path, f)
}

func (Native) Executable() (string, error) {
	return os.Executable() // TODO: not sure this works with the remote concepts
}

func (Native) Open(path string) (fs.File, error) {
	return os.Open(path)
}

func (Native) Create(path string, mode fs.FileMode) (WritableFile, error) {
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
}

func (Native) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(name, atime, mtime)
}

func (Native) Mkdir(path string, mode fs.FileMode) error {
	return os.Mkdir(path, mode)
}

func (Native) MkdirAll(path string, mode fs.FileMode) error {
	return os.MkdirAll(path, mode)
}

func (Native) Remove(path string) error {
	return os.Remove(path)
}

func (Native) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (Native) Abs(path string) (string, error) {
	return filepath.Abs(path)
}
