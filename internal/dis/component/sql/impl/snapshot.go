package impl

import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
	"go.tkw01536.de/pkglib/stream"
)

// SnapshotDB makes a snapshot of the given database into dest.
func (impl *Impl) SnapshotDB(ctx context.Context, progress io.Writer, dest io.Writer, database string) (e error) {
	return impl.do(ctx, progress, func(stack *dockerx.Stack) error {
		code := stack.Exec(
			ctx,
			stream.NewIOStream(dest, progress, nil),
			dockerx.ExecOptions{
				Service: "sql",
				Cmd:     impl.DumpExecutable,
				Args:    []string{"--databases", database},
			},
		)()
		if code != 0 {
			return fmt.Errorf("failed to execute dump: exit code %d", code)
		}
		return nil
	})
}
