package sql

//spellchecker:words context github wisski distillery internal component models dockerx pkglib errorsx stream
import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/stream"
)

// SnapshotDB makes a backup of the sql database into dest.
func (sql *SQL) SnapshotDB(ctx context.Context, progress io.Writer, dest io.Writer, database string) (e error) {
	stack, err := sql.OpenStack()
	if err != nil {
		return fmt.Errorf("failed to open stack: %w", err)
	}
	defer errorsx.Close(stack, &e, "stack")

	code := stack.Exec(
		ctx,
		stream.NewIOStream(dest, progress, nil),
		dockerx.ExecOptions{
			Service: "sql",
			Cmd:     dumpExecutable,
			Args:    []string{"--databases", database},
		},
	)()
	if code != 0 {
		return errSQLBackup
	}
	return nil
}
