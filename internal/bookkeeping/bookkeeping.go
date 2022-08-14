// Package bookkeeping implements reading and writing from the bookkeeping table
package bookkeeping

import (
	"database/sql/driver"
	"errors"
	"time"
)

// Instance is a WissKI Instance inside the bookkeeping table.
// It does not represent a running instance; it does not perform any validation.
type Instance struct {
	// NOTE: Modifying this struct requires a database migration.
	// This should nnever be done unless you know what you're doing.

	// Primary key for the instance
	Pk uint `gorm:"column:pk;primaryKey"`

	// time the instance was created
	Created time.Time `gorm:"column:created;autoCreateTime"`

	// slug of the system
	Slug string `gorm:"column:slug;not null;unique"`

	// email address of the system owner (if any)
	OwnerEmail string `gorm:"column:owner_email;type:varchar(320)"`

	// should we automatically enable updates for the system?
	AutoBlindUpdateEnabled SQLBit1 `gorm:"column:auto_blind_update_enabled;default:1"`

	// The filesystem path the system can be found under
	FilesystemBase string `gorm:"column:filesystem_base;not null"`

	// SQL Database credentials for the system
	SqlDatabase string `gorm:"column:sql_database;not null"`
	SqlUser     string `gorm:"column:sql_user;not null"`
	SqlPassword string `gorm:"column:sql_password;not null"`

	// GraphDB Repository
	GraphDBRepository string `gorm:"column:graphdb_repository;not null"`
	GraphDBUser       string `gorm:"column:graphdb_user;not null"`
	GraphDBPassword   string `gorm:"column:graphdb_password;not null"`
}

func (i Instance) IsBlindUpdateEnabled() bool {
	return bool(i.AutoBlindUpdateEnabled)
}

// SQLBit1 implements a boolean as a BIT(1)
type SQLBit1 bool

func (sb SQLBit1) Value() (driver.Value, error) {
	if sb {
		return []byte{1}, nil
	} else {
		return []byte{0}, nil
	}
}

var errBadBool = errors.New("SQLBit1: Database does not contain Bit(1)")

func (sb *SQLBit1) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok && len(bytes) == 1 {
		*sb = bytes[0] == 1
		return nil
	}
	return errBadBool
}
