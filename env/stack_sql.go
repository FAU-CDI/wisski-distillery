package env

import (
	"fmt"
	"io"
	"io/fs"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/bookkeeping"
	"github.com/FAU-CDI/wisski-distillery/internal/logging"
	"github.com/FAU-CDI/wisski-distillery/internal/sqle"
	"github.com/FAU-CDI/wisski-distillery/internal/stack"
	"github.com/FAU-CDI/wisski-distillery/internal/wait"
	"github.com/pkg/errors"
	"github.com/tkw1536/goprogram/stream"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SQLStack returns the docker stack that handles the sql database.
func (dis *Distillery) SQLStack() stack.Installable {
	return dis.asCoreStack("sql", stack.Installable{
		MakeDirsPerm: fs.ModeDir | fs.ModePerm,
		MakeDirs: []string{
			"data",
		},
	})
}

// SQLStackPath returns the path the SQLStack() lives at.
func (dis *Distillery) SQLStackPath() string {
	return dis.SQLStack().Dir
}

// sqlOpen opens a new sql connection to the provided database using the administrative credentials
func (env Distillery) sqlOpen(database string, config *gorm.Config) (*gorm.DB, error) {
	sql := mysql.Config{
		DSN:               fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", env.Config.MysqlAdminUser, env.Config.MysqlAdminPassword, "127.0.0.1:3306", database),
		DefaultStringSize: 256,
	}

	db, err := gorm.Open(mysql.New(sql), config)
	if err != nil {
		return db, err
	}

	gdb, err := db.DB()
	if err != nil {
		return db, err
	}
	gdb.SetMaxIdleConns(0)

	return db, nil
}

// sqlBkTable returns a gorm connection to the bookkeeping database.
func (dis *Distillery) sqlBkTable(silent bool) (*gorm.DB, error) {

	config := &gorm.Config{}
	if silent {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}

	// open the database
	db, err := dis.sqlOpen(dis.Config.DistilleryBookkeepingDatabase, config)
	if err != nil {
		return nil, err
	}

	// load the table
	table := db.Table(dis.Config.DistilleryBookkeepingTable)
	if table.Error != nil {
		return nil, err
	}

	return table, nil
}

var errSQLBackup = errors.New("SQLBackup: Mysqldump returned non-zero exit code")

// SQLBackup makes a backup of the sql database into dest.
func (dis *Distillery) SQLBackup(io stream.IOStream, dest io.Writer, database string) error {
	io = stream.NewIOStream(dest, io.Stderr, nil, 0)

	code, err := dis.SQLStack().Exec(io, "sql", "mysqldump", "--database", database)
	if err != nil {
		return err
	}
	if code != 0 {
		return errSQLBackup
	}
	return nil
}

// SQLShell executes a mysql shell inside the SQLStack.
func (dis *Distillery) SQLShell(io stream.IOStream, argv ...string) (int, error) {
	return dis.SQLStack().Exec(io, "sql", "mysql", argv...)
}

const waitSQLInterval = 1 * time.Second

// SQLWaitForShell waits for the sql database to be reachable via a docker-compose shell
func (dis *Distillery) SQLWaitForShell() error {
	n := stream.FromNil()
	return wait.Wait(func() bool {
		code, err := dis.SQLShell(n, "-e", "show databases;")
		return err == nil && code == 0
	}, waitSQLInterval, dis.Context())
}

// SQLWaitForConnection waits for the sql connection to be alive
func (dis *Distillery) SQLWaitForConnection() error {
	return wait.Wait(func() bool {
		_, err := dis.sqlBkTable(true)
		return err == nil
	}, waitSQLInterval, dis.Context())
}

var errInvalidDatabaseName = errors.New("SQLProvision: Invalid database name")

func (dis *Distillery) sqlRaw(query string, args ...interface{}) bool {
	sql := sqle.Format(query, args...)
	code, err := dis.SQLShell(stream.FromNil(), "-e", sql)
	return err == nil && code == 0
}

// SQLProvision provisions a new sql database and user
func (dis *Distillery) SQLProvision(name, user, password string) error {
	// wait for the database
	if err := dis.SQLWaitForShell(); err != nil {
		return err
	}

	// it's not a safe database name!
	if !sqle.IsSafeDatabaseName(name) {
		return errInvalidDatabaseName
	}

	// create the database and user!
	if !dis.sqlRaw("CREATE DATABASE `"+name+"`; CREATE USER ?@`%` IDENTIFIED BY ?; GRANT ALL PRIVILEGES ON `"+name+"`.* TO ?@`%`; FLUSH PRIVILEGES;", user, password, user) {
		return errors.New("SQLProvision: Failed to create user")
	}

	// and done!
	return nil
}

var errSQLPurgeUser = errors.New("unable to delete user")

// SQLPurgeUser deletes the specified user from the database
func (dis *Distillery) SQLPurgeUser(user string) error {
	if !dis.sqlRaw("DROP USER IF EXISTS ?@`%`; FLUSH PRIVILEGES; ", user) {
		return errSQLPurgeUser
	}

	return nil
}

var errSQLPurgeDB = errors.New("unable to drop database")

// SQLPurgeDatabase deletes the specified db from the database
func (dis *Distillery) SQLPurgeDatabase(db string) error {
	if !sqle.IsSafeDatabaseName(db) {
		return errSQLPurgeDB
	}
	if !dis.sqlRaw("DROP DATABASE IF EXISTS `" + db + "`") {
		return errSQLPurgeDB
	}
	return nil
}

var errSQLUnableToCreateUser = errors.New("unable to create administrative user")
var errSQLUnsafeDatabaseName = errors.New("Bookkeeping database has an unsafe name")
var errSQLUnableToCreate = errors.New("unable to create bookkeeping database")

// SQLBootstrap bootstraps the SQL database, and makes sure that the bookkeeping table is up-to-date
func (dis *Distillery) SQLBootstrap(io stream.IOStream) error {
	if err := dis.SQLWaitForShell(); err != nil {
		return err
	}

	// create the admin user
	logging.LogMessage(io, "Creating administrative user")
	{
		username := dis.Config.MysqlAdminUser
		password := dis.Config.MysqlAdminPassword
		if !dis.sqlRaw("CREATE USER IF NOT EXISTS ?@'%' IDENTIFIED BY ?; GRANT ALL PRIVILEGES ON *.* TO ?@`%` WITH GRANT OPTION; FLUSH PRIVILEGES;", username, password, username) {
			return errSQLUnableToCreateUser
		}
	}

	// create the admin user
	logging.LogMessage(io, "Creating sql database")
	{
		if !sqle.IsSafeDatabaseName(dis.Config.DistilleryBookkeepingDatabase) {
			return errSQLUnsafeDatabaseName
		}
		createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", dis.Config.DistilleryBookkeepingDatabase)
		if !dis.sqlRaw(createDBSQL) {
			return errSQLUnableToCreate
		}
	}

	// wait for the database to come up
	logging.LogMessage(io, "Waiting for database update to be complete")
	dis.SQLWaitForConnection()

	// open the database
	logging.LogMessage(io, "Migrating bookkeeping table")
	{
		db, err := dis.sqlBkTable(false)
		if err != nil {
			return fmt.Errorf("unable to access bookkeeping table: %s", err)
		}

		if err := db.AutoMigrate(&bookkeeping.Instance{}); err != nil {
			return fmt.Errorf("unable to migrate bookkeeping table: %s", err)
		}
	}

	return nil
}
