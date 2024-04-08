package sql

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
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

// unsafeWaitShell waits for a connection via the database shell to succeed
func (sql *SQL) unsafeWaitShell(ctx context.Context) (err error) {
	defer func() {
		// catch the errSQLNotFound
		r := recover()
		if r == nil {
			return
		}

		// other panic => keep panicking
		if r != errSQLNotFound {
			panic(r)
		}

		err = errSQLNotFound
	}()

	return timex.TickUntilFunc(func(time.Time) bool {
		code := sql.Shell(ctx, stream.FromNil(), "-e", "select 1;")

		// special case: executable was not found in the docker container.
		// so bail out immediately; as there is no hope of recovery.
		if code == 127 || code == 126 {
			panic(errSQLNotFound)
		}
		return code == 0
	}, ctx, sql.PollInterval)
}

// unsafeQuery shell executes a raw database query.
func (sql *SQL) unsafeQueryShell(ctx context.Context, query string) bool {
	code := sql.Shell(ctx, stream.FromNil(), "-e", query)
	return code == 0
}

var errSQLUnableToCreateUser = errors.New("unable to create administrative user")
var errSQLUnsafeDatabaseName = errors.New("distillery database has an unsafe name")
var errSQLUnableToMigrate = exit.Error{
	Message:  "unable to migrate %s table: %s",
	ExitCode: exit.ExitGeneric,
}

// Update initializes or updates the SQL database.
func (sql *SQL) Update(ctx context.Context, progress io.Writer) error {
	config := component.GetStill(sql).Config.SQL

	// unsafely create the admin user!
	{
		if err := sql.unsafeWaitShell(ctx); err != nil {
			return err
		}
		logging.LogMessage(progress, "Creating administrative user")
		{
			username := config.AdminUsername
			password := config.AdminPassword
			if err := sql.CreateSuperuser(ctx, username, password, true); err != nil {
				return errSQLUnableToCreateUser
			}
		}
	}

	// create the admin user
	logging.LogMessage(progress, "Creating sql database")
	{
		if !sqlx.IsSafeDatabaseLiteral(config.Database) {
			return errSQLUnsafeDatabaseName
		}
		createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", config.Database)
		if err := sql.Exec(createDBSQL); err != nil {
			return err
		}
	}

	// wait for the database to come up
	logging.LogMessage(progress, "Waiting for database update to be complete")
	if err := sql.WaitQueryTable(ctx); err != nil {
		return err
	}

	// migrate all of the tables!
	return logging.LogOperation(func() error {
		for _, table := range sql.dependencies.Tables {
			info := table.TableInfo()
			logging.LogMessage(progress, "migrating %q table", table.Name())
			db, err := sql.queryTable(ctx, false, info.Name)
			if err != nil {
				return errSQLUnableToMigrate.WithMessageF(table.Name, "unable to access table")
			}

			tp := reflect.New(info.Model).Interface()

			if err := db.AutoMigrate(tp); err != nil {
				return errSQLUnableToMigrate.WithMessageF(table.Name, err)
			}
		}
		return nil
	}, progress, "migrating database tables")
}
