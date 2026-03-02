package sqldelegator

import (
	"context"
	"io"
)

func (delegator *delegated) Snapshot(ctx context.Context, progress io.Writer, dest io.Writer) error {
	panic("not implemented")
}

func (delegator *delegated) Restore(ctx context.Context, progress io.Writer, src io.Reader) error {
	panic("not implemented")
}
