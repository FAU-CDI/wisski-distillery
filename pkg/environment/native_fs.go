package environment

import (
	"io/fs"
	"os"
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
