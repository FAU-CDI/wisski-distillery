package sqldelegator

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"regexp"

	"go.tkw01536.de/pkglib/stream"
)

func (delegated *delegated) Snapshot(ctx context.Context, progress io.Writer, dest io.Writer) error {
	return delegated.delegator.dependencies.SQL.SnapshotDB(ctx, progress, dest, delegated.instance.SqlDatabase)
}

func (delegated *delegated) Restore(ctx context.Context, reader io.Reader, io stream.IOStream) error {
	replacedFile := replaceSqlDatabaseName(reader, delegated.instance.SqlDatabase)
	defer replacedFile.Close()

	code := delegated.Shell(ctx, stream.NewIOStream(io.Stdout, io.Stderr, replacedFile))
	if code != 0 {
		return fmt.Errorf("failed to restore SQL contents: exit code %d", code)
	}
	return nil
}

var (
	reCreateDB = regexp.MustCompile(
		`(?i)^\s*CREATE\s+DATABASE\b.*?` + "`" + `([^` + "`" + `]+)` + "`",
	)
	reUseDB = regexp.MustCompile(
		`(?i)^\s*USE\s+(` + "`" + `([^` + "`" + `]+)` + "`" + `|([^\s;]+))`,
	)
)

func replaceSqlDatabaseName(reader io.Reader, newDB string) io.ReadCloser {
	// HACK HACK HACK: This restore code makes shit tons of assumptions about the SQL dump.
	// In particular that it was created by mysqldump -- and only one 'CREATE DATABASE' statement exists.

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		scanner := bufio.NewScanner(reader)

		// allow large lines (mysqldump can emit big INSERTs)
		scanner.Buffer(make([]byte, 64*1024), 10*1024*1024)

		for scanner.Scan() {
			line := scanner.Text()

			// CREATE DATABASE ... `oldDB`
			if m := reCreateDB.FindStringSubmatchIndex(line); m != nil {
				// group 1 = db name inside backticks
				dbStart, dbEnd := m[2], m[3]

				if line[dbStart:dbEnd] != "" {
					line = line[:dbStart] + newDB + line[dbEnd:]
				}
			} else if m := reUseDB.FindStringSubmatchIndex(line); m != nil {
				// USE `db`  -> group 2
				// USE db    -> group 3
				dbStart, dbEnd := -1, -1
				if m[4] != -1 { // backticked
					dbStart, dbEnd = m[4], m[5]
				} else if m[6] != -1 { // bare
					dbStart, dbEnd = m[6], m[7]
				}
				if dbStart != -1 && dbEnd != -1 && line[dbStart:dbEnd] != "" {
					line = line[:dbStart] + newDB + line[dbEnd:]
				}
			}

			if _, err := io.WriteString(pw, line+"\n"); err != nil {
				_ = pw.CloseWithError(err)
				return
			}
		}

		if err := scanner.Err(); err != nil {
			_ = pw.CloseWithError(err)
		}
	}()

	return pr
}
