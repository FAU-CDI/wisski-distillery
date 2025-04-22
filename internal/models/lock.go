//spellchecker:words models
package models

// LockTable is the name of the table the 'Metadatum' model is stored in.
const LockTable = "locks"

// Lock represents a log on WissKI Instances.
type Lock struct {
	Pk uint `gorm:"column:pk;primaryKey"`

	Slug string `gorm:"column:slug;not null"` // slug of instance
}
