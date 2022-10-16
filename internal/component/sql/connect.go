package sql

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/timex"
)

//
// ========== low-level connection ==========
//

// Exec executes a database-independent database query.
func (sql *SQL) Exec(query string, args ...interface{}) error {
	// connect to the server
	conn, err := sql.connect("")
	if err != nil {
		return err
	}

	// do the query!
	{
		_, err := conn.Exec(query, args...)
		if err != nil {
			return err
		}
		return nil
	}
}

// WaitExec waits for the query interface to be able to connect to the database
func (sql *SQL) WaitExec() error {
	return timex.TickUntilFunc(func(time.Time) bool {
		err := sql.Exec("select 1;")
		return err == nil
	}, sql.PollContext, sql.PollInterval)
}

//
// ========== connection via gorm ==========
//

// QueryTable returns a gorm.DB to connect to the provided distillery database table
func (sql *SQL) QueryTable(silent bool, table string) (*gorm.DB, error) {
	conn, err := sql.connect(sql.Config.DistilleryDatabase)
	if err != nil {
		return nil, err
	}

	// gorm configuration
	config := &gorm.Config{}
	if silent {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}

	// mysql connection
	cfg := mysql.Config{
		Conn: conn,

		DefaultStringSize: 256,
	}

	// open the gorm connection!
	db, err := gorm.Open(mysql.New(cfg), config)
	if err != nil {
		return nil, err
	}

	// set the table
	db = db.Table(table)

	// check that nothing went wrong
	if db.Error != nil {
		return nil, db.Error
	}
	return db, nil
}

// WaitQueryTable waits for a connection to succeed via QueryTable
func (sql *SQL) WaitQueryTable() error {
	// TODO: Establish a convention on when to wait for this!
	return timex.TickUntilFunc(func(time.Time) bool {
		_, err := sql.QueryTable(true, models.InstanceTable)
		return err == nil
	}, sql.PollContext, sql.PollInterval)
}

//
// ========== low-level database connection ==========
//

func (ssql *SQL) connect(database string) (*sql.DB, error) {
	conn, err := sql.Open("mysql", ssql.dsn(database))
	if err != nil {
		return nil, err
	}

	conn.SetMaxIdleConns(0)

	return conn, nil
}

// dsn returns a dsn fof connecting to the database
func (sql *SQL) dsn(database string) string {
	user := sql.Config.MysqlAdminUser
	pass := sql.Config.MysqlAdminPassword
	network := sql.network()
	server := sql.ServerURL

	return fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, network, server, database)
}

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
