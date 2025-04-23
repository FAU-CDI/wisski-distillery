package sql

//spellchecker:words context github wisski distillery internal component models pkglib stream
import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/tkw1536/pkglib/stream"
)

func (*SQL) SnapshotNeedsRunning() bool { return false }

func (*SQL) SnapshotName() string { return "sql" }

func (sql *SQL) Snapshot(wisski models.Instance, scontext *component.StagingContext) error {
	if err := scontext.AddDirectory(".", func(ctx context.Context) error {
		if err := scontext.AddFile(wisski.SqlDatabase+".sql", func(ctx context.Context, file io.Writer) error {
			if err := sql.SnapshotDB(ctx, scontext.Progress(), file, wisski.SqlDatabase); err != nil {
				return fmt.Errorf("failed to snapshot database: %w", err)
			}
			return nil
		}); err != nil {
			return fmt.Errorf("failed to add sql file: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to snapshot directory: %w", err)
	}
	return nil
}

// SnapshotDB makes a backup of the sql database into dest.
func (sql *SQL) SnapshotDB(ctx context.Context, progress io.Writer, dest io.Writer, database string) error {
	code := sql.Stack().Exec(ctx, stream.NewIOStream(dest, progress, nil), "sql", SQlDumpExecutable, "--databases", database)()
	if code != 0 {
		return errSQLBackup
	}
	return nil
}
