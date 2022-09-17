package sql

import (
	"errors"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
)

var errSQLBackup = errors.New("SQLBackup: Mysqldump returned non-zero exit code")

func (*SQL) BackupName() string {
	return "sql.sql"
}

// Backup makes a backup of all SQL databases into the path dest.
func (sql *SQL) Backup(context component.BackupContext) error {
	return context.AddFile("", func(file io.Writer) error {
		io := context.IO().Streams(file, nil, nil, 0).NonInteractive()
		code, err := sql.Stack().Exec(io, "sql", "mysqldump", "--all-databases")
		if err != nil {
			return err
		}
		if code != 0 {
			return errSQLBackup
		}
		return nil
	})

}
