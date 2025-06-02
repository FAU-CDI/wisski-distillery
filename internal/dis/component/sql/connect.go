package sql

//spellchecker:words context database time gorm driver mysql logger github wisski distillery internal component models pkglib timex
import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/tkw1536/pkglib/errorsx"
	"github.com/tkw1536/pkglib/timex"
)

//
// ========== low-level connection ==========
//

// Exec executes a database-independent database query.
func (sql *SQL) Exec(query string, args ...interface{}) (e error) {
	// connect to the server
	conn, err := sql.openConnection("")
	if err != nil {
		return err
	}
	defer errorsx.Close(conn, &e, "connection")

	// do the query!
	{
		_, err := conn.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
		return nil
	}
}

//
// ========== connection via gorm ==========
//

// QueryTableLegacy returns a gorm.DB to connect to the provided table of the given model.
// Deprecated: Use QueryTable instead.
func (sql *SQL) QueryTableLegacy(ctx context.Context, table component.Table) (*gorm.DB, error) {
	return sql.queryTable(ctx, queryTableOpts{
		table: table.TableInfo().Name,
	})
}

var errWrongGenericType = errors.New("wrong generic type for table")

// QueryTable returns a gorm.Context that has connected to the database.
func QueryTable[T any](ctx context.Context, sql *SQL, table component.Table) (gorm.Interface[T], error) {
	info := table.TableInfo()
	if got := reflect.TypeFor[T](); got != info.Model {
		return nil, fmt.Errorf("%w: got %v, expected %v", errWrongGenericType, got, info.Model)
	}

	db, err := sql.queryTable(ctx, queryTableOpts{table: info.Name})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to table: %w", err)
	}

	return gorm.G[T](db), nil
}

type queryTableOpts struct {
	silent bool
	table  string
}

// queryTable returns a gorm.DB to connect to the provided distillery database table.
func (sql *SQL) queryTable(ctx context.Context, opts queryTableOpts) (*gorm.DB, error) {
	conn, err := sql.connect()
	if err != nil {
		return nil, err
	}

	// gorm configuration
	config := &gorm.Config{
		Logger: newGormLogger(),
	}
	if opts.silent {
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
		return nil, fmt.Errorf("failed to connect with gorm: %w", err)
	}

	// set the table
	db = db.WithContext(ctx).Table(opts.table)

	// check that nothing went wrong
	if db.Error != nil {
		return nil, db.Error
	}
	return db, nil
}

// WaitQueryTable waits for a connection to succeed via QueryTable.
func (sql *SQL) WaitQueryTable(ctx context.Context) error {
	// TODO: Establish a convention on when to wait for this!
	if err := timex.TickUntilFunc(func(time.Time) bool {
		// TODO: Use a different table here
		_, err := sql.queryTable(ctx, queryTableOpts{table: models.InstanceTable})
		return err == nil
	}, ctx, sql.PollInterval); err != nil {
		return fmt.Errorf("failed to wait for connection: %w", err)
	}
	return nil
}

//
// ========== low-level database connection ==========
//

// connect establishes a connection to the sql database.
// the first succesfull connection is cached, and re-used automatically.
func (ssql *SQL) connect() (*sql.DB, error) {
	ssql.m.Lock()
	defer ssql.m.Unlock()

	if ssql.db != nil {
		return ssql.db, nil
	}

	database := component.GetStill(ssql).Config.SQL.Database
	dsn := ssql.dsn(database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// openConnection opens a new sql connection to the given database.
// The caller should manage the connection, and take care of its' lifecycle.
func (ssql *SQL) openConnection(database string) (*sql.DB, error) {
	db, err := sql.Open("mysql", ssql.dsn(database))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sql: %w", err)
	}
	return db, nil
}

// dsn returns a dsn fof connecting to the database.
func (sql *SQL) dsn(database string) string {
	config := component.GetStill(sql).Config.SQL
	user := config.AdminUsername
	pass := config.AdminPassword
	network := "tcp"
	server := sql.ServerURL

	return fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, network, server, database)
}
