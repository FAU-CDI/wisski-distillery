package sql

import (
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/tkw1536/goprogram/stream"
)

func (SQL) SnapshotNeedsRunning() bool { return false }

func (SQL) SnapshotName() string { return "sql" }

func (sql *SQL) Snapshot(wisski models.Instance, context component.StagingContext) error {
	return context.AddDirectory(".", func() error {
		return context.AddFile(wisski.SqlDatabase+".sql", func(file io.Writer) error {
			return sql.SnapshotDB(context.IO(), file, wisski.SqlDatabase)
		})
	})
}

// SnapshotDB makes a backup of the sql database into dest.
func (sql *SQL) SnapshotDB(io stream.IOStream, dest io.Writer, database string) error {
	io = io.Streams(dest, nil, nil, 0).NonInteractive()

	code, err := sql.Stack(sql.Environment).Exec(io, "sql", "mysqldump", "--databases", database)
	if err != nil {
		return err
	}
	if code != 0 {
		return errSQLBackup
	}
	return nil
}
