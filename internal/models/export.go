//spellchecker:words models
package models

//spellchecker:words time
import "time"

var _ Model = Export{}

// Export represents an entry in the export log.
type Export struct {
	Pk uint `gorm:"column:pk;primaryKey"`

	Slug    string    `gorm:"column:slug"`    // slug of instance
	Created time.Time `gorm:"column:created"` // time the backup was created

	Path   string `gorm:"column:path;not null"`   // path the export is stored at
	Packed bool   `gorm:"column:packed;not null"` // was the export packed, or was it staging only?
}

func (Export) TableName() string {
	return "snapshot"
}
