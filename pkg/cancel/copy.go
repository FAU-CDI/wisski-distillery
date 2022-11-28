package cancel

import (
	"context"
	"io"
	"time"
)

type SetDeadline interface {
	SetDeadline(t time.Time)
}

type SetReadDeadline interface {
	SetReadDeadline(t time.Time) error
}

type SetWriteDeadline interface {
	SetWriteDeadline(t time.Time) error
}

// Copy reads from src, and copies to dst.
//
// If the context is closed before src is closed, attempts to close the underlying reader and writer.
func Copy(ctx context.Context, dst io.Writer, src io.Reader) (written int64, err error) {

	// if the context has a deadline, the propanate that deadline to the underyling file.
	// this might cause the read call to not block.
	if deadline, ok := ctx.Deadline(); ok {
		var zero time.Time

		if file, ok := src.(SetReadDeadline); ok {
			file.SetReadDeadline(deadline)
			defer file.SetReadDeadline(zero)
		} else if file, ok := src.(SetDeadline); ok {
			file.SetDeadline(deadline)
			defer file.SetDeadline(zero)
		}

		if file, ok := dst.(SetWriteDeadline); ok {
			file.SetWriteDeadline(deadline)
			defer file.SetWriteDeadline(zero)
		} else if file, ok := dst.(SetDeadline); ok {
			file.SetDeadline(deadline)
			defer file.SetDeadline(zero)
		}
	}

	written, err, _ = WithContext2(ctx, func(start func()) (int64, error) {
		start()
		return io.Copy(dst, src)
	}, func() {
		if closer, ok := src.(io.Closer); ok {
			closer.Close()
		}
	})
	return written, err
}
