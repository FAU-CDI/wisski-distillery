//spellchecker:words models
package models

//spellchecker:words time
import (
	"fmt"
	"time"

	"github.com/tkw1536/pkglib/fsx"
)

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

// Exists checks if the given export exists on disk.
func (e Export) Exists() (bool, error) {
	if e.Path == "" {
		return false, nil
	}

	exists, err := fsx.Exists(e.Path)
	if err != nil {
		return false, fmt.Errorf("unable to check existence: %w", err)
	}
	return exists, nil
}
