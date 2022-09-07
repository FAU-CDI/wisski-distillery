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

// SQLComponent represents the 'sql' layer belonging to a distillery
type SQLComponent struct {
	PollInterval time.Duration // Duration to wait for during wait

	dis *Distillery
}

// SSH returns the SSHComponent belonging to this distillery
func (dis *Distillery) SQL() SQLComponent {
	return SQLComponent{
		PollInterval: time.Second,

		dis: dis,
	}
}

func (SQLComponent) Name() string {
	return "sql"
}

// Stack returns the docker stack that handles the sql database.
func (sql SQLComponent) Stack() stack.Installable {
	return sql.dis.makeComponentStack(sql, stack.Installable{
		MakeDirsPerm: fs.ModeDir | fs.ModePerm,
		MakeDirs: []string{
			"data",
		},
	})
}

// SQLStackPath returns the path the SQLStack() lives at.
func (sql SQLComponent) Path() string {
	return sql.Stack().Dir
}

// sqlOpen opens a new sql connection to the provided database using the administrative credentials
func (sql SQLComponent) openDatabase(database string, config *gorm.Config) (*gorm.DB, error) {
	cfg := mysql.Config{
		DSN:               fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", sql.dis.Config.MysqlAdminUser, sql.dis.Config.MysqlAdminPassword, "127.0.0.1:3306", database),
		DefaultStringSize: 256,
	}

	db, err := gorm.Open(mysql.New(cfg), config)
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

// OpenBookkeeping opens a connection to the bookkeeping database
func (sql SQLComponent) OpenBookkeeping(silent bool) (*gorm.DB, error) {

	config := &gorm.Config{}
	if silent {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}

	// open the database
	db, err := sql.openDatabase(sql.dis.Config.DistilleryBookkeepingDatabase, config)
	if err != nil {
		return nil, err
	}

	// load the table
	table := db.Table(sql.dis.Config.DistilleryBookkeepingTable)
	if table.Error != nil {
		return nil, err
	}

	return table, nil
}

var errSQLBackup = errors.New("SQLBackup: Mysqldump returned non-zero exit code")

// Backup makes a backup of the sql database into dest.
func (sql SQLComponent) Backup(io stream.IOStream, dest io.Writer, database string) error {
	io = stream.NewIOStream(dest, io.Stderr, nil, 0)

	code, err := sql.Stack().Exec(io, "sql", "mysqldump", "--databases", database)
	if err != nil {
		return err
	}
	if code != 0 {
		return errSQLBackup
	}
	return nil
}

// BackupAll makes a backup of all sql databases
func (sql SQLComponent) BackupAll(io stream.IOStream, dest io.Writer) error {
	io = stream.NewIOStream(dest, io.Stderr, nil, 0)

	code, err := sql.Stack().Exec(io, "sql", "mysqldump", "--all-databases")
	if err != nil {
		return err
	}
	if code != 0 {
		return errSQLBackup
	}
	return nil
}

// OpenShell executes a mysql shell command
func (sql SQLComponent) OpenShell(io stream.IOStream, argv ...string) (int, error) {
	return sql.Stack().Exec(io, "sql", "mysql", argv...)
}

// WaitShell waits for the sql database to be reachable via a docker-compose shell
func (sql SQLComponent) WaitShell() error {
	n := stream.FromNil()
	return wait.Wait(func() bool {
		code, err := sql.OpenShell(n, "-e", "show databases;")
		return err == nil && code == 0
	}, sql.PollInterval, sql.dis.Context())
}

// Wait waits for a connection to the bookkeeping table to suceed
func (sql SQLComponent) Wait() error {
	return wait.Wait(func() bool {
		_, err := sql.OpenBookkeeping(true)
		return err == nil
	}, sql.PollInterval, sql.dis.Context())
}

var errInvalidDatabaseName = errors.New("SQLProvision: Invalid database name")

func (sql SQLComponent) Query(query string, args ...interface{}) bool {
	raw := sqle.Format(query, args...)
	code, err := sql.OpenShell(stream.FromNil(), "-e", raw)
	return err == nil && code == 0
}

// SQLProvision provisions a new sql database and user
func (sql SQLComponent) Provision(name, user, password string) error {
	// wait for the database
	if err := sql.WaitShell(); err != nil {
		return err
	}

	// it's not a safe database name!
	if !sqle.IsSafeDatabaseName(name) {
		return errInvalidDatabaseName
	}

	// create the database and user!
	if !sql.Query("CREATE DATABASE `"+name+"`; CREATE USER ?@`%` IDENTIFIED BY ?; GRANT ALL PRIVILEGES ON `"+name+"`.* TO ?@`%`; FLUSH PRIVILEGES;", user, password, user) {
		return errors.New("SQLProvision: Failed to create user")
	}

	// and done!
	return nil
}

var errSQLPurgeUser = errors.New("unable to delete user")

// SQLPurgeUser deletes the specified user from the database
func (sql SQLComponent) PurgeUser(user string) error {
	if !sql.Query("DROP USER IF EXISTS ?@`%`; FLUSH PRIVILEGES; ", user) {
		return errSQLPurgeUser
	}

	return nil
}

var errSQLPurgeDB = errors.New("unable to drop database")

// SQLPurgeDatabase deletes the specified db from the database
func (sql SQLComponent) PurgeDatabase(db string) error {
	if !sqle.IsSafeDatabaseName(db) {
		return errSQLPurgeDB
	}
	if !sql.Query("DROP DATABASE IF EXISTS `" + db + "`") {
		return errSQLPurgeDB
	}
	return nil
}

var errSQLUnableToCreateUser = errors.New("unable to create administrative user")
var errSQLUnsafeDatabaseName = errors.New("Bookkeeping database has an unsafe name")
var errSQLUnableToCreate = errors.New("unable to create bookkeeping database")

// Bootstrap bootstraps the SQL database, and makes sure that the bookkeeping table is up-to-date
func (sql SQLComponent) Bootstrap(io stream.IOStream) error {
	if err := sql.WaitShell(); err != nil {
		return err
	}

	// create the admin user
	logging.LogMessage(io, "Creating administrative user")
	{
		username := sql.dis.Config.MysqlAdminUser
		password := sql.dis.Config.MysqlAdminPassword
		if !sql.Query("CREATE USER IF NOT EXISTS ?@'%' IDENTIFIED BY ?; GRANT ALL PRIVILEGES ON *.* TO ?@`%` WITH GRANT OPTION; FLUSH PRIVILEGES;", username, password, username) {
			return errSQLUnableToCreateUser
		}
	}

	// create the admin user
	logging.LogMessage(io, "Creating sql database")
	{
		if !sqle.IsSafeDatabaseName(sql.dis.Config.DistilleryBookkeepingDatabase) {
			return errSQLUnsafeDatabaseName
		}
		createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", sql.dis.Config.DistilleryBookkeepingDatabase)
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
