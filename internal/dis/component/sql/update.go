package sql

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"

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
	return sql.Stack().Exec(ctx, io, "sql", "mysql", argv...)()
}

// unsafeWaitShell waits for a connection via the database shell to succeed
func (sql *SQL) unsafeWaitShell(ctx context.Context) error {
	n := stream.FromNil()
	return timex.TickUntilFunc(func(time.Time) bool {
		code := sql.Shell(ctx, n, "-e", "select 1;")
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

	// unsafely create the admin user!
	{
		if err := sql.unsafeWaitShell(ctx); err != nil {
			return err
		}
		logging.LogMessage(progress, "Creating administrative user")
		{
			username := sql.Config.SQL.AdminUsername
			password := sql.Config.SQL.AdminPassword
			if err := sql.CreateSuperuser(ctx, username, password, true); err != nil {
				return errSQLUnableToCreateUser
			}
		}
	}

	// create the admin user
	logging.LogMessage(progress, "Creating sql database")
	{
		if !sqlx.IsSafeDatabaseLiteral(sql.Config.SQL.Database) {
			return errSQLUnsafeDatabaseName
		}
		createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", sql.Config.SQL.Database)
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
		for _, table := range sql.Dependencies.Tables {
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
