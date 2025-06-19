package sql

import (
	"context"
	"errors"
	"fmt"
)

// CreateOpts are options for [CreateOpts].
type CreateOpts struct {
	Name        string // Name of the database to create
	AllowExists bool   // Don't error if the database already exists

	CreateUser bool // Create an appropriate database user
	Username   string
	Password   string
}

var (
	errCreateOptsMissingName = errors.New("CreateOpts: missing name")
	errCreateOptsCreateUser  = errors.New("CreateOpts: username and password must be given if and only if createUser is true")
)

func (cd CreateOpts) Validate() error {
	if cd.Name == "" {
		return errCreateOptsMissingName
	}

	if cd.CreateUser != (cd.Username != "" && cd.Password != "") {
		return errCreateOptsCreateUser
	}

	return nil
}

// CreateDatabase creates a new database with the given name.
// If the user name is not the empty string, it then generates a new user and grants access to this database.
//
// Provision internally waits for the database to become available.
func (sql *SQL) CreateDatabase(ctx context.Context, opts CreateOpts) error {
	// NOTE(twiesing): We shouldn't use string concat to build sql queries.
	// But the driver doesn't support using query params for these queries.
	// Apparently it's a "feature", see https://github.com/go-sql-driver/mysql/issues/398#issuecomment-169951763.

	if err := opts.Validate(); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	var (
		nameQuoted = quoteBacktick(opts.Name)
	)

	IfNotExists := ""
	if opts.AllowExists {
		IfNotExists = "IF NOT EXISTS"
	}

	queries := []string{
		"CREATE DATABASE " + IfNotExists + " " + nameQuoted + ";",
	}
	if opts.CreateUser {
		var (
			userQuoted = quoteSingle(opts.Username)
			passQuoted = quoteSingle(opts.Password)
		)

		queries = append(queries,
			"CREATE USER "+userQuoted+"@`%` IDENTIFIED BY "+passQuoted+";",
			"GRANT ALL PRIVILEGES ON "+nameQuoted+".* TO "+userQuoted+"@`%`;",
			"FLUSH PRIVILEGES;",
		)
	}

	return sql.directQuery(ctx, queries...)
}

// Drops the given database if it exists.
func (sql *SQL) DropDatabase(ctx context.Context, db string) error {
	var (
		dbQuoted = quoteBacktick(db)
	)

	return sql.directQuery(ctx, "DROP DATABASE IF EXISTS "+dbQuoted+";")
}

// DropUser drops the given user if it exists.
func (sql *SQL) DropUser(ctx context.Context, user string) error {
	var (
		userQuoted = quoteSingle(user)
	)

	return sql.directQuery(
		ctx,
		"DROP USER IF EXISTS "+userQuoted+"@'%';",
		"FLUSH PRIVILEGES;",
	)
}
