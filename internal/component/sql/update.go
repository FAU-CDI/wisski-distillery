package sql

import (
	"errors"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/FAU-CDI/wisski-distillery/pkg/sqle"
	"github.com/FAU-CDI/wisski-distillery/pkg/wait"
	"github.com/tkw1536/goprogram/exit"
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
		n.EPrintf("[SQL.unsafeWaitShell]: %d %s\n", code, err)
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
var errSQLUnableToMigrate = exit.Error{
	Message:  "unable to migrate %s table: %s",
	ExitCode: exit.ExitGeneric,
}

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
		if err := sql.Exec(createDBSQL); err != nil {
			return err
		}
	}

	// wait for the database to come up
	logging.LogMessage(io, "Waiting for database update to be complete")
	sql.WaitQueryTable()

	tables := []struct {
		name  string
		model any
		table string
	}{
		{
			"instance",
			&models.Instance{},
			models.InstanceTable,
		},
		{
			"metadata",
			&models.Metadatum{},
			models.MetadataTable,
		},
		{
			"snapshot",
			&models.Export{},
			models.ExportTable,
		},
		{
			"lock",
			&models.Lock{},
			models.LockTable,
		},
	}

	// migrate all of the tables!
	return logging.LogOperation(func() error {
		for _, table := range tables {
			logging.LogMessage(io, "migrating %q table", table.name)
			db, err := sql.QueryTable(false, table.table)
			if err != nil {
				return errSQLUnableToMigrate.WithMessageF(table.name, "unable to access table")
			}

			if err := db.AutoMigrate(table.model); err != nil {
				return errSQLUnableToMigrate.WithMessageF(table.name, err)
			}
		}
		return nil
	}, io, "migrating database tables")
}
