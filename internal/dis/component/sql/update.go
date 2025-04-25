package sql

//spellchecker:words context errors reflect time github wisski distillery internal component logging goprogram exit pkglib sqlx stream timex
import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/pkglib/sqlx"
	"github.com/tkw1536/pkglib/stream"
	"github.com/tkw1536/pkglib/timex"
)

// Shell runs a mysql shell with the provided databases.
//
// NOTE(twiesing): This command should not be used to connect to the database or execute queries except in known situations.
func (sql *SQL) Shell(ctx context.Context, io stream.IOStream, argv ...string) int {
	return sql.Stack().Exec(ctx, io, "sql", "mariadb", argv...)()
}

var errSQLNotFound = errors.New("internal error: unsafeWaitShell: sql client not found")

// unsafeWaitShell waits for a connection via the database shell to succeed.
func (sql *SQL) unsafeWaitShell(ctx context.Context) (err error) {
	defer func() {
		// catch the errSQLNotFound
		r := recover()
		if r == nil {
			return
		}

		// if we simply didn't find the sql, don't panic!
		if e, ok := r.(error); ok && errors.Is(e, errSQLNotFound) {
			err = e
			return
		}
	}()

	if err := timex.TickUntilFunc(func(time.Time) bool {
		code := sql.Shell(ctx, stream.FromNil(), "-e", "select 1;")

		// special case: executable was not found in the docker container.
		// so bail out immediately; as there is no hope of recovery.
		if code == 127 || code == 126 {
			panic(errSQLNotFound)
		}
		return code == 0
	}, ctx, sql.PollInterval); err != nil {
		return fmt.Errorf("failed to wait for sql: %w", err)
	}
	return nil
}

// unsafeQuery shell executes a raw database query.
func (sql *SQL) unsafeQueryShell(ctx context.Context, query string) bool {
	code := sql.Shell(ctx, stream.FromNil(), "-e", query)
	return code == 0
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
		if err := sql.unsafeWaitShell(ctx); err != nil {
			return err
		}
		if _, err := logging.LogMessage(progress, "Creating administrative user"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		{
			username := config.AdminUsername
			password := config.AdminPassword
			if err := sql.CreateSuperuser(ctx, username, password, true); err != nil {
				return errSQLUnableToCreateUser
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
	} //  shouldn't abort cause logging failed
	if err := sql.WaitQueryTable(ctx); err != nil {
		return fmt.Errorf("failed to wait for database: %w", err)
	}

	// migrate all of the tables!
	if err := logging.LogOperation(func() error {
		for _, table := range sql.dependencies.Tables {
			info := table.TableInfo()
			if _, err := logging.LogMessage(progress, "migrating %q table", table.Name()); err != nil {
				return fmt.Errorf("failed to log message: %w", err)
			}
			db, err := sql.queryTable(ctx, false, info.Name)
			if err != nil {
				return fmt.Errorf("failed to access table %q for migration: %w", table.Name(), err)
			}

			tp := reflect.New(info.Model).Interface()

			if err := db.AutoMigrate(tp); err != nil {
				return fmt.Errorf("failed auto migration for table %q: %w", table.Name(), err)
			}
		}
		return nil
	}, progress, "migrating database tables"); err != nil {
		return fmt.Errorf("failed to migrate database tables: %w", err)
	}
	return nil
}
