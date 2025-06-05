//spellchecker:words models
package models

var _ Model = Lock{}

// Lock represents a lock on WissKI Instances.
type Lock struct {
	Pk uint `gorm:"column:pk;primaryKey"`

	Slug string `gorm:"column:slug;not null;unique"` // slug of instance
}

func (Lock) TableName() string {
	return "locks"
}
