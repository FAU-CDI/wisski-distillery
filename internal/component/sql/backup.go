package sql

import (
	"errors"
	"os"

	"github.com/tkw1536/goprogram/stream"
)

var errSQLBackup = errors.New("SQLBackup: Mysqldump returned non-zero exit code")

func (*SQL) BackupName() string {
	return "sql.sql"
}

// Backup makes a backup of all SQL databases into the path dest.
func (sql *SQL) Backup(io stream.IOStream, dest string) error {
	// open the file, TODO: Outsource this to context
	writer, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer writer.Close()

	// run sqldump
	io = io.Streams(writer, nil, nil, 0).NonInteractive()
	code, err := sql.Stack().Exec(io, "sql", "mysqldump", "--all-databases")
	if err != nil {
		return err
	}
	if code != 0 {
		return errSQLBackup
	}
	return nil
}
