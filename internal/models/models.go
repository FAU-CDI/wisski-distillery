// Package contains all database models
//
//spellchecker:words models
package models

// Model represents an abitrary database model
type Model interface {
	//
	TableName() string
}
