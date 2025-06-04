package sql

//spellchecker:words context database time gorm driver mysql logger github wisski distillery internal component models pkglib timex
import (
	"context"
	databaseSQL "database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/tkw1536/pkglib/errorsx"
	"github.com/tkw1536/pkglib/timex"
)

//
// ========== low-level connection ==========
//

// Exec executes a database-independent database query.
func (sql *SQL) Exec(query string, args ...interface{}) (e error) {
	// connect to the server
	conn, err := sql.openSQL("")
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

// Wait waits for a connection to the sql database to succeed.
// The database doesn't have to exist.
func (sql *SQL) Wait(ctx context.Context) (err error) {
	if err := timex.TickUntilFunc(func(time.Time) bool {
		conn, err := sql.openSQL("")
		if err != nil {
			return false
		}
		defer conn.Close()

		if _, err := conn.QueryContext(ctx, "select 1;"); err != nil {
			return false
		}
		return true
	}, ctx, sql.PollInterval); err != nil {
		return fmt.Errorf("failed to wait for sql: %w", err)
	}
	return nil
}

//
// ========== connection via gorm ==========
//

// OpenTable opens a *gorm.DB connection to the given table.
// Should use [OpenInterface] where possible.
//
// TODO: Migrate everything here!
func (sql *SQL) OpenTable(ctx context.Context, table component.Table) (*gorm.DB, error) {
	db, err := sql.connectGorm(ctx)
	if err != nil {
		return nil, err
	}

	name := table.TableInfo().Name()

	db = db.Table(name)
	if db.Error != nil {
		return nil, fmt.Errorf("failed to open connection to table %q: %w", name, db.Error)
	}

	return db, nil
}

var errWrongGenericType = errors.New("wrong generic type for table")

// OpenInterface opens a [gorm.Interface] to the given sql and table interface.
// The generic parameter T must correspond to the [component.Table]'s TableInfo.
func OpenInterface[T any](ctx context.Context, sql *SQL, table component.Table) (gorm.Interface[T], error) {
	info := table.TableInfo()
	if got := reflect.TypeFor[T](); got != reflect.TypeOf(info.Model) {
		return nil, fmt.Errorf("%w: got %v, expected %v", errWrongGenericType, got, info.Model)
	}

	db, err := sql.connectGorm(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to table: %w", err)
	}

	return gorm.G[T](db), nil
}

// connectGorm returns a gorm.DB to connect to the provided distillery database table.
func (sql *SQL) connectGorm(ctx context.Context) (*gorm.DB, error) {
	conn, err := sql.connectSQL()
	if err != nil {
		return nil, err
	}

	// gorm configuration
	config := &gorm.Config{
		Logger: newGormLogger().LogMode(logger.Info),
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

	db = db.WithContext(ctx)

	// check that nothing went wrong
	if db.Error != nil {
		return nil, db.Error
	}
	return db, nil
}

//
// ========== low-level database connection ==========
//

// connectSQL establishes a connection to the sql database.
// the first succesfull connection is cached, and re-used automatically.
func (sql *SQL) connectSQL() (*databaseSQL.DB, error) {
	// slow path: not yet loaded
	if !sql.dbOpen.Load() {
		return sql.connectSQLSlow()
	}

	// fast path: already open!
	return sql.db, nil
}

func (sql *SQL) connectSQLSlow() (*databaseSQL.DB, error) {
	sql.m.Lock()
	defer sql.m.Unlock()

	// someone else grabbed the lock before us
	// and already opened the db!
	if sql.db != nil {
		return sql.db, nil
	}

	db, err := sql.openSQL(
		component.GetStill(sql).Config.SQL.Database,
	)
	if err != nil {
		return nil, err
	}

	// store the db for future invocations
	// and tell everyone to take the fast path
	sql.db = db
	sql.dbOpen.Store(true)

	return db, nil
}

// openSQL opens a new sql connection to the given database.
// The caller should manage the connection, and take care of its' lifecycle.
func (sql *SQL) openSQL(database string) (*databaseSQL.DB, error) {
	config := component.GetStill(sql).Config.SQL

	var (
		user    = config.AdminUsername
		pass    = config.AdminPassword
		network = "tcp"
		server  = sql.ServerURL
	)

	db, err := databaseSQL.Open(
		"mysql",
		fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, network, server, database),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database %q: %w", database, err)
	}
	return db, nil
}
