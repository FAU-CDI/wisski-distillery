package models

import (
	"database/sql/driver"
	"errors"
	"time"
)

// InstanceTable is the name of the table the 'Instance' model is stored in.
const InstanceTable = "distillery"

// Instance is a WissKI Instance stored inside the sql database.
//
// It does not represent a running instance; it does not perform any validation.
type Instance struct {
	// NOTE: Modifying this struct requires a database migration.
	// This should never be done unless you know what you're doing.

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

	// DockerBaseImage is the php base image to use
	DockerBaseImage string `gorm:"column:docker_base;not_null"`

	// SQL Database credentials for the system
	SqlDatabase string `gorm:"column:sql_database;not null"`
	SqlUsername string `gorm:"column:sql_user;not null"`
	SqlPassword string `gorm:"column:sql_password;not null"`

	// GraphDB Repository
	GraphDBRepository string `gorm:"column:graphdb_repository;not null"`
	GraphDBUsername   string `gorm:"column:graphdb_user;not null"`
	GraphDBPassword   string `gorm:"column:graphdb_password;not null"`
}

const (
	PHP8         = "8.0"
	PHP8_IMAGE   = "docker.io/library/php:8.0-apache-bullseye"
	PHP8_1       = "8.1"
	PHP8_1_IMAGE = "docker.io/library/php:8.1-apache-bullseye"
)

var errUnknownPHPVersion = errors.New("unknown php version")

// GetBaseImage returns the php base image to use
func GetBaseImage(php string) (string, error) {
	switch php {
	case "":
		return PHP8_IMAGE, nil
	case PHP8:
		return PHP8_IMAGE, nil
	case PHP8_1:
		return PHP8_1_IMAGE, nil
	default:
		return "", errUnknownPHPVersion
	}
}

func (i Instance) GetDockerBaseImage() string {
	if i.DockerBaseImage == "" {
		return PHP8_IMAGE
	}
	return i.DockerBaseImage
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

var errBadBool = errors.New("`SQLBit1': database does not contain `Bit(1)'")

func (sb *SQLBit1) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok && len(bytes) == 1 {
		*sb = bytes[0] == 1
		return nil
	}
	return errBadBool
}
