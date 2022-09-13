package countwriter

import (
	"io"
)

// CountWriter wraps an io.Writer, see [NewCountWriter].
//
// It is intended to be used to count different writes to an underlying writer.
// Once an error occurs, no more writes are passed through, and the underlying error is returned instead.
// This means that in practice, calls to write can be continued and are ignored silently.
//
// The underlying sum of bytes written and error can be seen using [Sum].
type CountWriter struct {
	w io.Writer

	n   int
	err error
}

// NewCountWriter creates a new [CountWriter] that delegates to w.
func NewCountWriter(w io.Writer) *CountWriter {
	return &CountWriter{w: w}
}

// write performs the write operation w on this writer.
func (cw *CountWriter) write(w func() (int, error)) (int, error) {
	// if there was an error, return it and don't do a write
	if cw.err != nil {
		return 0, cw.err
	}

	// call the writer
	n, err := w()

	// update the underling state
	cw.n += n
	cw.err = err

	// and return
	return n, err
}

// Write implements [io.Writer]
func (cw *CountWriter) Write(p []byte) (int, error) {
	return cw.write(func() (int, error) {
		return cw.w.Write(p)
	})
}

// WriteString implements [io.WriteString].
// See [Write].
func (cw *CountWriter) WriteString(s string) (int, error) {
	return cw.write(func() (int, error) {
		return io.WriteString(cw.w, s)
	})
}

// Sum returns the state, that is the total number of bytes written and any error
func (cw *CountWriter) Sum() (int, error) {
	return cw.n, cw.err
}
