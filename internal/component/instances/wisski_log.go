package instances

import (
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/tkw1536/goprogram/lib/collection"
)

// ExportLogFor retrieves (and prunes) the ExportLog.
// Slug determines if entries for Backups (empty slug)
// or a specific Instance (non-empty slug) are returned.
func (instances *Instances) ExportLogFor(slug string) (exports []models.Export, err error) {
	exports, err = instances.ExportLog()
	if err != nil {
		return nil, err
	}

	return collection.Filter(exports, func(s models.Export) bool {
		return s.Slug == slug
	}), nil
}

// ExportLog retrieves (and prunes) all entries in the snapshot log.
func (instances *Instances) ExportLog() ([]models.Export, error) {
	// query the table!
	table, err := instances.SQL.QueryTable(false, models.ExportTable)
	if err != nil {
		return nil, err
	}

	// find all the exports
	var exports []models.Export
	res := table.Find(&exports)
	if res.Error != nil {
		return nil, res.Error
	}

	// partition out the exports that have been deleted!
	parts := collection.Partition(exports, func(s models.Export) bool {
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
func (wisski *WissKI) Snapshots() (snapshots []models.Export, err error) {
	return wisski.instances.ExportLogFor(wisski.Slug)
}

// AddToExportLog adds the provided export to the log.
func (instances *Instances) AddToExportLog(export models.Export) error {
	// find the table
	table, err := instances.SQL.QueryTable(false, models.ExportTable)
	if err != nil {
		return err
	}

	// and save it!
	res := table.Create(&export)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
