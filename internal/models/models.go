// Package contains all database models
//
//spellchecker:words models
package models

//spellchecker:words gorm schema
import "gorm.io/gorm/schema"

// Model represents an abitrary database model.
type Model interface {
	schema.Tabler
}
