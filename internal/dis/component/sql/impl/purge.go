package impl

import (
	"context"
	"io"
	"strings"
)

// Purge purges the given database and user.
// If they do not exist, they are not purged.
func (impl *Impl) Purge(ctx context.Context, progress io.Writer, database string, user string) error {
	return impl.queries(
		ctx,
		progress,
		"DROP DATABASE IF EXISTS "+quoteBacktick(database)+";",
		"DROP USER IF EXISTS "+quoteSingle(user)+"@'%';",
		"FLUSH PRIVILEGES;",
	)
}

// quoteBacktick quotes value using backticks.
func quoteBacktick(value string) string {
	var builder strings.Builder
	_, _ = builder.WriteRune('`')

	for _, r := range value {
		if r == '`' {
			_, _ = builder.WriteRune('`')
		}

		_, _ = builder.WriteRune(r)
	}

	_, _ = builder.WriteRune('`')
	return builder.String()
}

// quoteSingle quotes values for use in a mariadb single quoted string.
func quoteSingle(value string) string {
	var builder strings.Builder
	_, _ = builder.WriteRune('\'')

	for _, r := range value {
		if r == '\'' || r == '\\' {
			builder.WriteRune('\\')
		}
		_, _ = builder.WriteRune(r)
	}

	_, _ = builder.WriteRune('\'')
	return builder.String()
}
