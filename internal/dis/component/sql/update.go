package sql

//spellchecker:words context errors reflect time github wisski distillery internal component execx logging pkglib sqlx stream timex
import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
	"github.com/FAU-CDI/wisski-distillery/pkg/execx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/pkglib/errorsx"
	"github.com/tkw1536/pkglib/sqlx"
	"github.com/tkw1536/pkglib/stream"
)

// Shell runs a mysql shell with the provided databases.
//
// NOTE(twiesing): This command should not be used to connect to the database or execute queries except in known situations.
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
			Cmd:     "mariadb",
			Args:    argv,
		},
	)()
}

// unsafeQuoteBacktick quotes value using a backtick.
func unsafeQuoteBacktick(value string) string {
	return "`" + strings.ReplaceAll(value, "`", "``") + "`"
}

// unsafeQuoteSingle quotes value using a single tickmark
func unsafeQuoteSingle(value string) string {
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

// directQuery opens a new connection to the database and executes the given queries in order.
// Once the queries have been executed, the connection is closed.
func (sql *SQL) directQuery(ctx context.Context, queries ...string) (e error) {
	conn, err := sql.openSQL("")
	if err != nil {
		return fmt.Errorf("failed to establish connection: %w", err)
	}
	defer errorsx.Close(conn, &e, "connection")

	for _, query := range queries {
		if _, err := conn.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}

var (
	errSQLUnableToCreateUser = errors.New("unable to create administrative user")
	errSQLUnsafeDatabaseName = errors.New("distillery database has an unsafe name")
)

// Update initializes or updates the SQL database.
func (sql *SQL) Update(ctx context.Context, progress io.Writer) error {
	config := component.GetStill(sql).Config.SQL

	// unsafely create the admin user!
	{
		if err := sql.Wait(ctx); err != nil {
			return err
		}
		if _, err := logging.LogMessage(progress, "Creating administrative user"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		{
			username := config.AdminUsername
			password := config.AdminPassword
			if err := sql.CreateSuperuser(ctx, username, password, true); err != nil {
				return fmt.Errorf("%w: %w", errSQLUnableToCreateUser, err)
			}
		}
	}

	// create the admin user
	if _, err := logging.LogMessage(progress, "Creating sql database"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	} //  shouldn't abort cause logging failed
	{
		if !sqlx.IsSafeDatabaseLiteral(config.Database) {
			return errSQLUnsafeDatabaseName
		}
		createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", config.Database)
		if err := sql.Exec(createDBSQL); err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
	}

	// wait for the database to come up
	if _, err := logging.LogMessage(progress, "Waiting for database update to be complete"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}

	if err := sql.Wait(ctx); err != nil {
		return fmt.Errorf("failed to wait for database: %w", err)
	}

	// migrate all of the tables!
	if err := logging.LogOperation(func() error {
		for _, table := range sql.dependencies.Tables {
			info := table.TableInfo()
			table := info.Name()

			if _, err := logging.LogMessage(progress, "migrating %q table", table); err != nil {
				return fmt.Errorf("failed to log message: %w", err)
			}

			db, err := sql.connectGorm(ctx)
			if err != nil {
				return fmt.Errorf("failed to connect to table %q: %w", table, err)
			}

			if err := db.AutoMigrate(info.Model); err != nil {
				return fmt.Errorf("failed auto migration for table %q: %w", table, err)
			}
		}
		return nil
	}, progress, "migrating database tables"); err != nil {
		return fmt.Errorf("failed to migrate database tables: %w", err)
	}
	return nil
}
