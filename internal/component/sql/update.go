package sql

import (
	"errors"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/bookkeeping"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/FAU-CDI/wisski-distillery/pkg/sqle"
	"github.com/tkw1536/goprogram/stream"
)

var errSQLUnableToCreateUser = errors.New("unable to create administrative user")
var errSQLUnsafeDatabaseName = errors.New("bookkeeping database has an unsafe name")
var errSQLUnableToCreate = errors.New("unable to create bookkeeping database")

// Update initializes or updates the SQL database.
func (sql *SQL) Update(io stream.IOStream) error {
	if err := sql.WaitShell(); err != nil {
		return err
	}

	// create the admin user
	logging.LogMessage(io, "Creating administrative user")
	{
		username := sql.Config.MysqlAdminUser
		password := sql.Config.MysqlAdminPassword
		if !sql.Query("CREATE USER IF NOT EXISTS ?@'%' IDENTIFIED BY ?; GRANT ALL PRIVILEGES ON *.* TO ?@`%` WITH GRANT OPTION; FLUSH PRIVILEGES;", username, password, username) {
			return errSQLUnableToCreateUser
		}
	}

	// create the admin user
	logging.LogMessage(io, "Creating sql database")
	{
		if !sqle.IsSafeDatabaseName(sql.Config.DistilleryBookkeepingDatabase) {
			return errSQLUnsafeDatabaseName
		}
		createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", sql.Config.DistilleryBookkeepingDatabase)
		if !sql.Query(createDBSQL) {
			return errSQLUnableToCreate
		}
	}

	// wait for the database to come up
	logging.LogMessage(io, "Waiting for database update to be complete")
	sql.Wait()

	// open the database
	logging.LogMessage(io, "Migrating bookkeeping table")
	{
		db, err := sql.OpenBookkeeping(false)
		if err != nil {
			return fmt.Errorf("unable to access bookkeeping table: %s", err)
		}

		if err := db.AutoMigrate(&bookkeeping.Instance{}); err != nil {
			return fmt.Errorf("unable to migrate bookkeeping table: %s", err)
		}
	}

	return nil
}
