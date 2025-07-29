package sql

//spellchecker:words context strings github wisski distillery dockerx execx pkglib stream
import (
	"context"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
	"github.com/FAU-CDI/wisski-distillery/pkg/execx"
	"go.tkw01536.de/pkglib/stream"
)

// Shell directly executes a mysql command inside the container.
// This command should be used with caution.
func (sql *SQL) Shell(ctx context.Context, io stream.IOStream, argv ...string) int {
	stack, err := sql.OpenStack()
	if err != nil {
		return execx.CommandError
	}
	defer func() {
		_ = stack.Close()
	}()

	return stack.Exec(
		ctx, io,
		dockerx.ExecOptions{
			Service: "sql",
			Cmd:     queryExecutable,
			Args:    argv,
		},
	)()
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
