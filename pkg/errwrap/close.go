package errwrap

import (
	"errors"
	"fmt"
	"io"
	"runtime"
)

type failedToCloseError struct {
	file string
	line int

	name string
	err  error
}

func (ftc failedToCloseError) Error() string {
	if ftc.file == "" {
		return fmt.Sprintf("failed to close %s: %s", ftc.name, ftc.err)
	}
	return fmt.Sprintf("failed to close %s near %s:%d: %s", ftc.name, ftc.file, ftc.line, ftc.err)
}

func (ftc failedToCloseError) Unwrap() error {
	return ftc.err
}

// Close closes the given closer and updates retval if closing failed.
// desc is used as the description of closer.
//
// Close is intended to be defered:
//
//	func stuff() (e error) {
//		f, err := os.Open(...)
//		if err != nil { /* ... */ }
//		defer errwrap.Close(f, "file", &e)
//		/* ... */
//	}
func Close(closer io.Closer, desc string, retval *error) {
	if retval == nil {
		panic("Close: nil retval should be replaced by a plain .Close() call")
	}

	err := closer.Close()
	if err == nil {
		return
	}

	// get stack trace info
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = ""
		line = 0
	}

	err = failedToCloseError{file: file, line: line, name: desc, err: err}

	if *retval == nil {
		*retval = err
	} else {
		*retval = errors.Join(*retval, err)
	}
}
