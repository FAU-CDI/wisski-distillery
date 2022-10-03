package instances

import (
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/slicesx"
)

// SnapshotLogFor retrieves (and prunes) the SnapshotLog for the provided slug.
// An empty slug returns the log of backups.
func (instances *Instances) SnapshotLogFor(slug string) (snapshots []models.Snapshot, err error) {
	snapshots, err = instances.SnapshotLog()
	if err != nil {
		return nil, err
	}

	return slicesx.Filter(snapshots, func(s models.Snapshot) bool {
		return s.Slug == slug
	}), nil
}

// SnapshotLog retrieves (and prunes) all entries in the snapshot log.
func (instances *Instances) SnapshotLog() ([]models.Snapshot, error) {
	// query the table!
	table, err := instances.SQL.QueryTable(false, models.SnapshotTable)
	if err != nil {
		return nil, err
	}

	// find all the snapshots
	var snapshots []models.Snapshot
	res := table.Find(&snapshots)
	if res.Error != nil {
		return nil, res.Error
	}

	// partition out the snapshots that have been deleted!
	parts := slicesx.Partition(snapshots, func(s models.Snapshot) bool {
		_, err := instances.Core.Environment.Stat(s.Path)
		return !environment.IsNotExist(err)
	})

	// go and delete them!
	if len(parts[false]) > 0 {
		if err := table.Delete(parts[false]).Error; err != nil {
			return nil, err
		}
	}

	// return the ones that still exist
	return parts[true], nil
}

// Snapshots returns the list of snapshots of this WissKI
func (wisski *WissKI) Snapshots() (snapshots []models.Snapshot, err error) {
	return wisski.instances.SnapshotLogFor(wisski.Slug)
}

// AddSnapshotLog adds a log entry for the provided entry
func (instances *Instances) AddSnapshotLog(snapshot models.Snapshot) error {
	// find the table
	table, err := instances.SQL.QueryTable(false, models.SnapshotTable)
	if err != nil {
		return err
	}

	// and save it!
	res := table.Create(&snapshot)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
