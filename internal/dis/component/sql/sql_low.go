package sql

import (
	"context"
	databaseSQL "database/sql"
	"errors"
	"fmt"
	"time"

	mysql "github.com/go-sql-driver/mysql"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/timex"
)

// Wait waits for the connection to the distillery-specific database to succeed.
func (sql *SQL) Wait(ctx context.Context) (err error) {
	if err := timex.TickUntilFunc(func(time.Time) bool {
		_, err = sql.connectSQL(ctx)
		return err == nil
	}, ctx, sql.PollInterval); err != nil {
		return fmt.Errorf("failed to wait for sql: %w", err)
	}
	return nil
}

// Query executes a database-independent query.
func (sql *SQL) Query(ctx context.Context, query string, args ...interface{}) (e error) {
	// connect to the server
	conn, err := sql.openSQL("")
	if err != nil {
		return err
	}
	defer errorsx.Close(conn, &e, "connection")

	// do the query!
	{
		_, err := conn.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
		return nil
	}
}

// directQuery waits to establish a new connection to the database, and then executes the given queries in order.
// Once the queries have been executed, the connection is closed.
func (sql *SQL) directQuery(ctx context.Context, queries ...string) (e error) {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("failed to query: %w", err)
	}

loop:
	conn, err := sql.openSQL("")
	if err != nil {
		select {
		case <-ctx.Done():
			return fmt.Errorf("failed to wait for sql: %w", ctx.Err())
		case <-time.After(sql.PollInterval):
			goto loop
		}
	}

	for _, query := range queries {
		if _, err := conn.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}

// connectSQL establishes a connection to the sql database.
// The context is used to check connection validity, and not attached to the connection permanently.
func (sql *SQL) connectSQL(ctx context.Context) (*databaseSQL.DB, error) {
	sql.m.Lock()
	defer sql.m.Unlock()

	if sql.db != nil {
		pingErr := sql.db.PingContext(ctx)
		if pingErr == nil {
			return sql.db, nil
		}

		wdlog.Of(ctx).Warn(
			"connection was closed, establishing new connection",
			"error", pingErr,
		)

		// connection went away (we can't ping it anymore)
		sql.db = nil
	}

	db, err := sql.connectSQLImpl(ctx)
	if err != nil {
		return nil, err
	}

	// store the db for future invocations
	// and tell everyone to take the fast path
	sql.db = db

	return db, nil
}

// actual implemention of connecting to the sql.
// Use [sql.connectSQL].
func (sql *SQL) connectSQLImpl(ctx context.Context) (*databaseSQL.DB, error) {
	db, err := sql.openSQL(
		component.GetStill(sql).Config.SQL.Database,
	)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to establish connection: %w", err)
	}

	return db, nil
}

var errOpenSQLConfig = errors.New("invalid sql configuration")

const (
	maxOpenConns      = 10
	maxIdleLifetime   = time.Minute
	maxActiveLifetime = time.Hour
)

// openSQL opens a new sql connection to the given database.
// The caller should manage the connection, and take care of its' lifecycle.
func (sql *SQL) openSQL(database string) (*databaseSQL.DB, error) {
	sqlConfig := component.GetStill(sql).Config.SQL

	config := mysql.NewConfig()
	{
		config.Net = "tcp"
		config.Addr = sql.ServerURL

		config.User = sqlConfig.AdminUsername
		config.Passwd = sqlConfig.AdminPassword

		config.DBName = database

		config.ParseTime = true
		config.Loc = time.UTC

		if err := config.Apply(mysql.Charset("utf8", "")); err != nil {
			return nil, fmt.Errorf("%w: %w", errOpenSQLConfig, err)
		}
	}

	db, err := databaseSQL.Open(
		"mysql", config.FormatDSN(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database %q: %w", database, err)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxIdleTime(maxIdleLifetime)
	db.SetConnMaxLifetime(maxActiveLifetime)

	return db, nil
}
