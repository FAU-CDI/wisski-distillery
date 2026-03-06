package impl

import (
	"bufio"
	"io"
	"regexp"
)

var (
	reCreateDB = regexp.MustCompile(
		`(?i)^\s*CREATE\s+DATABASE\b.*?` + "`" + `([^` + "`" + `]+)` + "`",
	)
	reUseDB = regexp.MustCompile(
		`(?i)^\s*USE\s+(` + "`" + `([^` + "`" + `]+)` + "`" + `|([^\s;]+))`,
	)
)

// ReplaceSqlDatabaseName replaces the database name in a SQL dump in the given reader using the given function.
// If there are no database names to replace, the reader is returned unchanged.
func ReplaceSqlDatabaseName(reader io.Reader, replaceFunc func(string) string) io.ReadCloser {
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
					line = line[:dbStart] + replaceFunc(line[dbStart:dbEnd]) + line[dbEnd:]
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
					line = line[:dbStart] + replaceFunc(line[dbStart:dbEnd]) + line[dbEnd:]
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
