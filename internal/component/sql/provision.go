package sql

import (
	"errors"

	"github.com/FAU-CDI/wisski-distillery/pkg/sqle"
)

// SQLProvision provisions a new sql database and user
func (sql *SQL) Provision(name, user, password string) error {
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
func (sql *SQL) PurgeUser(user string) error {
	if !sql.Query("DROP USER IF EXISTS ?@`%`; FLUSH PRIVILEGES; ", user) {
		return errSQLPurgeUser
	}

	return nil
}

var errSQLPurgeDB = errors.New("unable to drop database")

// SQLPurgeDatabase deletes the specified db from the database
func (sql *SQL) PurgeDatabase(db string) error {
	if !sqle.IsSafeDatabaseName(db) {
		return errSQLPurgeDB
	}
	if !sql.Query("DROP DATABASE IF EXISTS `" + db + "`") {
		return errSQLPurgeDB
	}
	return nil
}
