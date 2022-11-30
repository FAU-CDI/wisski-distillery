package sql

import (
	"context"
	"errors"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/tkw1536/goprogram/stream"
)

var errSQLBackup = errors.New("SQLBackup: Mysqldump returned non-zero exit code")

func (*SQL) BackupName() string {
	return "sql.sql"
}

// Backup makes a backup of all SQL databases into the path dest.
func (sql *SQL) Backup(scontext component.StagingContext) error {
	return scontext.AddFile("", func(ctx context.Context, file io.Writer) error {
		code, err := sql.Stack(sql.Environment).Exec(ctx, stream.NonInteractive(scontext.Progress()), "sql", "mysqldump", "--all-databases")
		if err != nil {
			return err
		}
		if code != 0 {
			return errSQLBackup
		}
		return nil
	})

}
