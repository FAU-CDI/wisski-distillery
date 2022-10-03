package models

import "time"

// SnapshotTable is the name of the table the [SnapshotLog] model is stored in
const SnapshotTable = "snapshot"

// Snapshot represents an entry in the snapshot log
type Snapshot struct {
	Pk uint `gorm:"column:pk;primaryKey"`

	Slug    string    `gorm:"column:slug"`    // slug of instance
	Created time.Time `gorm:"column:created"` // time the backup was created

	Path   string `gorm:"column:path;not null"`   // path the backup is stored at
	Packed bool   `gorm:"column:packed;not null"` // was the backup packed, or was it staging only?

}
