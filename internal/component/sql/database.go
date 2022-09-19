package sql

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync/atomic"

	mysqldriver "github.com/go-sql-driver/mysql"

	"github.com/FAU-CDI/wisski-distillery/pkg/sqle"
	"github.com/FAU-CDI/wisski-distillery/pkg/wait"
	"github.com/tkw1536/goprogram/stream"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var proxyNameCounter uint64

// network returns the network to use to connect to the database
func (sql *SQL) network() string {
	return sql.lazyNetwork.Get(func() (name string) {
		network := "tcp"

		// register a new DialContext function to use the environment.
		// this seems like a bit of a hack, but it works for now.
		name = fmt.Sprintf("sql-network-%d", atomic.AddUint64(&proxyNameCounter, 1))
		mysqldriver.RegisterDialContext(name, func(ctx context.Context, addr string) (net.Conn, error) {
			return sql.Core.Environment.DialContext(ctx, network, addr)
		})
		return
	})
}

// sqlOpen opens a new sql connection to the provided database using the administrative credentials
func (sql *SQL) openDatabase(database string, config *gorm.Config) (*gorm.DB, error) {
	cfg := mysql.Config{
		DriverName:        "mysql",
		DSN:               fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8&parseTime=True&loc=Local", sql.Config.MysqlAdminUser, sql.Config.MysqlAdminPassword, sql.network(), sql.ServerURL, database),
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
func (sql *SQL) OpenBookkeeping(silent bool) (*gorm.DB, error) {

	config := &gorm.Config{}
	if silent {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}

	// open the database
	db, err := sql.openDatabase(sql.Config.DistilleryBookkeepingDatabase, config)
	if err != nil {
		return nil, err
	}

	// load the table
	table := db.Table(sql.Config.DistilleryBookkeepingTable)
	if table.Error != nil {
		return nil, err
	}

	return table, nil
}

// Shell runs a mysql shell command.
func (sql *SQL) Shell(io stream.IOStream, argv ...string) (int, error) {
	return sql.Stack(sql.Environment).Exec(io, "sql", "mysql", argv...)
}

// WaitShell waits for the sql database to be reachable via shell
func (sql *SQL) WaitShell() error {
	n := stream.FromNil()
	return wait.Wait(func() bool {
		code, err := sql.Shell(n, "-e", "show databases;")
		return err == nil && code == 0
	}, sql.PollInterval, sql.PollContext)
}

// Wait waits for a connection to the bookkeeping table to suceed
func (sql *SQL) Wait() error {
	return wait.Wait(func() bool {
		_, err := sql.OpenBookkeeping(true)
		return err == nil
	}, sql.PollInterval, sql.PollContext)
}

var errInvalidDatabaseName = errors.New("SQLProvision: Invalid database name")

// Query performs a raw database query
func (sql *SQL) Query(query string, args ...interface{}) bool {
	raw := sqle.Format(query, args...)
	code, err := sql.Shell(stream.FromNil(), "-e", raw)
	return err == nil && code == 0
}
