package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/tkw1536/pkglib/timex"
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

//
// ========== connection via gorm ==========
//

// QueryTable returns a gorm.DB to connect to the provided table of the given model
func (sql *SQL) QueryTable(ctx context.Context, table component.Table) (*gorm.DB, error) {
	return sql.queryTable(ctx, false, table.TableInfo().Name)
}

// queryTable returns a gorm.DB to connect to the provided distillery database table
func (sql *SQL) queryTable(ctx context.Context, silent bool, table string) (*gorm.DB, error) {
	conn, err := sql.connect(component.GetStill(sql).Config.SQL.Database)
	if err != nil {
		return nil, err
	}

	// gorm configuration
	config := &gorm.Config{
		Logger: newGormLogger(),
	}
	if silent {
		config.Logger = config.Logger.LogMode(logger.Silent)
	} else {
		config.Logger = config.Logger.LogMode(logger.Info)
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
	db = db.WithContext(ctx).Table(table)

	// check that nothing went wrong
	if db.Error != nil {
		return nil, db.Error
	}
	return db, nil
}

// WaitQueryTable waits for a connection to succeed via QueryTable
func (sql *SQL) WaitQueryTable(ctx context.Context) error {
	// TODO: Establish a convention on when to wait for this!
	return timex.TickUntilFunc(func(time.Time) bool {
		// TODO: Use a different table here
		_, err := sql.queryTable(ctx, true, models.InstanceTable)
		return err == nil
	}, ctx, sql.PollInterval)
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
	config := component.GetStill(sql).Config.SQL
	user := config.AdminUsername
	pass := config.AdminPassword
	network := "tcp"
	server := sql.ServerURL

	return fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, network, server, database)
}
