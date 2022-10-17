package wisski

import "github.com/FAU-CDI/wisski-distillery/internal/models"

// Snapshots returns the list of snapshots of this WissKI
func (wisski *WissKI) Snapshots() (snapshots []models.Export, err error) {
	return wisski.SnapshotsLog.For(wisski.Slug)
}
