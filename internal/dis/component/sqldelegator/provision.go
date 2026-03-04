package sqldelegator

import (
	"context"
	"io"
	"strings"

	"go.tkw01536.de/pkglib/stream"
)

func (delegated *delegated) SQLUrl() string {
	return "mysql://" + delegated.instance.SqlUsername + ":" + delegated.instance.SqlPassword + "@sql/" + delegated.instance.SqlDatabase
}

func (delegated *delegated) Provision(ctx context.Context) error {
	return delegated.Impl.CreateDatabase(ctx, stream.Null, CreateOpts{
		Name: delegated.instance.SqlDatabase,

		CreateUser: true,
		Username:   delegated.instance.SqlUsername,
		Password:   delegated.instance.SqlPassword,
	})
}

func (delegated *delegated) Purge(ctx context.Context) error {
	return delegated.Impl.Purge(ctx, stream.Null, delegated.instance.SqlDatabase, delegated.instance.SqlUsername)
}

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
