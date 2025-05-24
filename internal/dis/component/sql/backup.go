package sql

//spellchecker:words context errors github wisski distillery internal component pkglib errorsx stream
import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/tkw1536/pkglib/errorsx"
	"github.com/tkw1536/pkglib/stream"
)

var errSQLBackup = errors.New("`SQLBackup': mysqldump returned non-zero exit code")

func (*SQL) BackupName() string {
	return "sql.sql"
}

// Backup makes a backup of all SQL databases into the path dest.
func (sql *SQL) Backup(scontext *component.StagingContext) error {
	if err := scontext.AddFile("", func(ctx context.Context, file io.Writer) (e error) {
		stack, err := sql.OpenStack()
		if err != nil {
			return fmt.Errorf("failed to open stack: %w", err)
		}
		defer errorsx.Close(stack, &e, "stack")

		code := stack.Exec(ctx, stream.NewIOStream(file, scontext.Progress(), nil), "sql", SQlDumpExecutable, "--all-databases")()
		if code != 0 {
			return errSQLBackup
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to add to context: %w", err)
	}
	return nil
}
