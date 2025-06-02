package sql

//spellchecker:words context errors github wisski distillery internal models pkglib errorsx sqlx
import (
	"context"
	"errors"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/tkw1536/pkglib/errorsx"
	"github.com/tkw1536/pkglib/sqlx"
)

var errProvisionInvalidDatabaseParams = errors.New("`Provision': invalid parameters")
var errProvisionInvalidGrant = errors.New("`Provision': grant failed")

// Provision provisions sql-specific resource for the given instance.
func (sql *SQL) Provision(ctx context.Context, instance models.Instance, domain string) error {
	return sql.CreateDatabase(ctx, instance.SqlDatabase, instance.SqlUsername, instance.SqlPassword)
}

// Purge purges sql-specific resources for the given instance.
func (sql *SQL) Purge(ctx context.Context, instance models.Instance, domain string) error {
	return errorsx.Combine(
		sql.PurgeDatabase(instance.SqlDatabase),
		sql.PurgeUser(ctx, instance.SqlUsername),
	)
}

// CreateDatabase creates a new database with the given name.
// It then generates a new user, with the name 'user' and the password 'password', that is then granted access to this database.
//
// Provision internally waits for the database to become available.
func (sql *SQL) CreateDatabase(ctx context.Context, name, user, password string) error {
	// NOTE(twiesing): We shouldn't use string concat to build sql queries.
	// But the driver doesn't support using query params for these queries.
	// Apparently it's a "feature", see https://github.com/go-sql-driver/mysql/issues/398#issuecomment-169951763.

	// quick and dirty check to make sure that all the names won't sql inject.
	if !sqlx.IsSafeDatabaseLiteral(name) || !sqlx.IsSafeDatabaseSingleQuote(user) || !sqlx.IsSafeDatabaseSingleQuote(password) {
		return errProvisionInvalidDatabaseParams
	}

	if err := sql.waitDatabase(ctx); err != nil {
		return err
	}

	if err := sql.directQuery(ctx,
		"CREATE DATABASE `"+name+"`;",
		"CREATE USER '"+user+"'@'%' IDENTIFIED BY '"+password+"';",
		"GRANT ALL PRIVILEGES ON `"+name+"`.* TO `"+user+"`@`%`;",
		"FLUSH PRIVILEGES;",
	); err != nil {
		return fmt.Errorf("%w: %w", errProvisionInvalidGrant, err)
	}
	return nil
}

var errCreateSuperuserGrant = errors.New("`CreateSuperUser': grant failed")

// CreateSuperuser createsa new user, with the name 'user' and the password 'password'.
// It then grants this user superuser status in the database.
//
// CreateSuperuser internally waits for the database to become available.
func (sql *SQL) CreateSuperuser(ctx context.Context, user, password string, allowExisting bool) error {
	// NOTE(twiesing): This function unsafely uses the shell directly to create a superuser.
	// This is for two reasons:
	// (1) this is used during bootstraping
	// (2) The underlying driver doesn't support "GRANT ALL PRIVILEGES"
	// See also [sql.Provision].

	if !sqlx.IsSafeDatabaseSingleQuote(user) || !sqlx.IsSafeDatabaseSingleQuote(password) {
		return errProvisionInvalidDatabaseParams
	}

	if err := sql.waitDatabase(ctx); err != nil {
		return err
	}

	var IfNotExists string
	if allowExisting {
		IfNotExists = "IF NOT EXISTS"
	}

	if err := sql.directQuery(ctx,
		"CREATE USER "+IfNotExists+" '"+user+"'@'%' IDENTIFIED BY '"+password+"';",
		"GRANT ALL PRIVILEGES ON *.* TO '"+user+"'@'%' WITH GRANT OPTION;",
		"FLUSH PRIVILEGES;",
	); err != nil {
		return fmt.Errorf("%w: %w", errCreateSuperuserGrant, err)
	}
	return nil
}

var errPurgeUser = errors.New("`PurgeUser': failed to drop user")

// SQLPurgeUser deletes the specified user from the database.
func (sql *SQL) PurgeUser(ctx context.Context, user string) error {
	if !sqlx.IsSafeDatabaseSingleQuote(user) {
		return errPurgeUser
	}

	if err := sql.directQuery(ctx,
		"DROP USER IF EXISTS '"+user+"'@'%';",
		"FLUSH PRIVILEGES;",
	); err != nil {
		return fmt.Errorf("%w: %w", errPurgeUser, err)
	}

	return nil
}

var errSQLPurgeDB = errors.New("unable to drop database: unsafe database name")

// SQLPurgeDatabase deletes the specified db from the database.
func (sql *SQL) PurgeDatabase(db string) error {
	if !sqlx.IsSafeDatabaseLiteral(db) {
		return errSQLPurgeDB
	}
	return sql.Exec("DROP DATABASE IF EXISTS `" + db + "`")
}
