package impl

import (
	"context"
	"errors"
	"fmt"
	"io"
)

// CreateOpts are options for [CreateOpts].
type CreateOpts struct {
	Name        string // Name of the database to create
	AllowExists bool   // Don't error if the database or user already exist

	CreateUser bool // Create an appropriate database user
	Superuser  bool // makes the user a superuser
	Username   string
	Password   string
}

var (
	errCreateOptsMissingName = errors.New("missing name")
	errCreateOptsCreateUser  = errors.New("username and password must be given if and only if createUser is true")
	errCreateOptsSuperuser   = errors.New("can't create a superuser without creating a user")
)

func (cd CreateOpts) Validate() error {
	if cd.Name == "" {
		return errCreateOptsMissingName
	}

	if cd.CreateUser != (cd.Username != "" && cd.Password != "") {
		return errCreateOptsCreateUser
	}
	if cd.Superuser && !cd.CreateUser {
		return errCreateOptsSuperuser
	}

	return nil
}

// CreateDatabase creates a new database with the given name.
// If the user name is not the empty string, it then generates a new user and grants access to this database.
//
// Provision internally waits for the database to become available.
func (impl *Impl) CreateDatabase(ctx context.Context, progress io.Writer, opts CreateOpts) error {
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
			privHost   = nameQuoted
		)

		userIfNotExists := ""
		if opts.AllowExists {
			userIfNotExists = "IF NOT EXISTS "
		}

		grantWithOption := ""
		if opts.Superuser {
			privHost = "*"
			grantWithOption = " WITH GRANT OPTION"
		}

		queries = append(queries,
			"CREATE USER "+userIfNotExists+userQuoted+"@`%` IDENTIFIED BY "+passQuoted+";",
			"GRANT ALL PRIVILEGES ON "+privHost+".* TO "+userQuoted+"@`%`"+grantWithOption+";",
			"FLUSH PRIVILEGES;",
		)
	}

	return impl.queries(ctx, progress, queries...)
}
