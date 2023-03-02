// Package fsx provides additional file system functionality.
//
// All functions in this package ignore the umask.
// As such it is not safe to use otherwise equivalent functions provided by the standard go library concurrently with this package.
// Users should take care that no other code in their application uses these functions.
package fsx

import (
	"io/fs"
	"sync"
	"syscall"
)

// DefaultFilePerm should be used by callers to use a consistent file mode for new files.
const DefaultFilePerm fs.FileMode = 0666

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
