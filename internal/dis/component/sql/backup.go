package sql

//spellchecker:words context errors github wisski distillery internal component pkglib stream
import (
	"context"
	"errors"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/tkw1536/pkglib/stream"
)

var errSQLBackup = errors.New("`SQLBackup': mysqldump returned non-zero exit code")

func (*SQL) BackupName() string {
	return "sql.sql"
}

// Backup makes a backup of all SQL databases into the path dest.
func (sql *SQL) Backup(scontext *component.StagingContext) error {
	return scontext.AddFile("", func(ctx context.Context, file io.Writer) error {
		code := sql.Stack().Exec(ctx, stream.NewIOStream(file, scontext.Progress(), nil), "sql", SQlDumpExecutable, "--all-databases")()
		if code != 0 {
			return errSQLBackup
		}
		return nil
	})

}
