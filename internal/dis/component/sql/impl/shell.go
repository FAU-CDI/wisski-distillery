package impl

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
	"github.com/FAU-CDI/wisski-distillery/pkg/execx"
	"go.tkw01536.de/pkglib/stream"
)

// Shell executes a mysql shell command inside the sql implementation.
func (impl *Impl) Shell(ctx context.Context, io stream.IOStream, argv ...string) int {
	code := execx.CommandError
	err := impl.do(ctx, stream.Null, func(stack *dockerx.Stack) error {
		code = stack.Exec(ctx, io, dockerx.ExecOptions{
			Service: impl.Service,
			Cmd:     impl.QueryExecutable,
			Args:    argv,
		})()
		return nil
	})
	if err != nil {
		code = execx.CommandError
	}
	return code
}
