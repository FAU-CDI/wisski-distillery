package sqldelegator

import (
	"context"

	"go.tkw01536.de/pkglib/stream"
)

func (delegated *delegated) Shell(ctx context.Context, io stream.IOStream, argv ...string) int {
	return delegated.delegator.dependencies.SQL.Shell(ctx, io)
}
