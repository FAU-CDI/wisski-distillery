package sql

import (
	"errors"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/FAU-CDI/wisski-distillery/pkg/sqle"
	"github.com/FAU-CDI/wisski-distillery/pkg/wait"
	"github.com/tkw1536/goprogram/stream"
)

// Shell runs a mysql shell with the provided databases.
//
// NOTE(twiesing): This command should not be used to connect to the database or execute queries except in known situations.
func (sql *SQL) Shell(io stream.IOStream, argv ...string) (int, error) {
	return sql.Stack(sql.Environment).Exec(io, "sql", "mysql", argv...)
}

// unsafeWaitShell waits for a connection via the database shell to succeed
func (sql *SQL) unsafeWaitShell() error {
	n := stream.FromNil()
	return wait.Wait(func() bool {
		code, err := sql.Shell(n, "-e", "select 1;")
		// log.Printf("[unsafeWaitShell] %d %s\n", code, err) // debug
		return err == nil && code == 0
	}, sql.PollInterval, sql.PollContext)
}

// unsafeQuery shell executes a raw database query.
func (sql *SQL) unsafeQueryShell(query string) bool {
	code, err := sql.Shell(stream.FromNil(), "-e", query)
	return err == nil && code == 0
}

var errSQLUnableToCreateUser = errors.New("unable to create administrative user")
var errSQLUnsafeDatabaseName = errors.New("distillery database has an unsafe name")

// Update initializes or updates the SQL database.
func (sql *SQL) Update(io stream.IOStream) error {

	// unsafely create the admin user!
	{
		if err := sql.unsafeWaitShell(); err != nil {
			return err
		}
		logging.LogMessage(io, "Creating administrative user")
		{
			username := sql.Config.MysqlAdminUser
			password := sql.Config.MysqlAdminPassword
			if err := sql.CreateSuperuser(username, password, true); err != nil {
				return errSQLUnableToCreateUser
			}
		}
	}

	// create the admin user
	logging.LogMessage(io, "Creating sql database")
	{
		if !sqle.IsSafeDatabaseLiteral(sql.Config.DistilleryDatabase) {
			return errSQLUnsafeDatabaseName
		}
		createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", sql.Config.DistilleryDatabase)
		if err := sql.Query(createDBSQL); err != nil {
			return err
		}
	}

	// wait for the database to come up
	logging.LogMessage(io, "Waiting for database update to be complete")
	sql.WaitQueryTable()

	// open the database
	logging.LogMessage(io, "Migrating instances table")
	{
		db, err := sql.QueryTable(false, models.InstanceTable)
		if err != nil {
			return fmt.Errorf("unable to access bookkeeping table: %s", err)
		}

		if err := db.AutoMigrate(&models.Instance{}); err != nil {
			return fmt.Errorf("unable to migrate bookkeeping table: %s", err)
		}
	}

	return nil
}
