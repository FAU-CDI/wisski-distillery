package sql

import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/tkw1536/goprogram/stream"
)

func (*SQL) SnapshotNeedsRunning() bool { return false }

func (*SQL) SnapshotName() string { return "sql" }

func (sql *SQL) Snapshot(wisski models.Instance, scontext component.StagingContext) error {
	return scontext.AddDirectory(".", func(ctx context.Context) error {
		return scontext.AddFile(wisski.SqlDatabase+".sql", func(ctx context.Context, file io.Writer) error {
			return sql.SnapshotDB(ctx, scontext.Progress(), file, wisski.SqlDatabase)
		})
	})
}

// SnapshotDB makes a backup of the sql database into dest.
func (sql *SQL) SnapshotDB(ctx context.Context, progress io.Writer, dest io.Writer, database string) error {
	code := sql.Stack(sql.Environment).Exec(ctx, stream.NonInteractive(progress), "sql", "mysqldump", "--databases", database)()
	if code != 0 {
		return errSQLBackup
	}
	return nil
}
