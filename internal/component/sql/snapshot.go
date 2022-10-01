package sql

import (
	"io"

	"github.com/tkw1536/goprogram/stream"
)

// SnapshotDB makes a backup of the sql database into dest.
func (sql SQL) SnapshotDB(io stream.IOStream, dest io.Writer, database string) error {
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
