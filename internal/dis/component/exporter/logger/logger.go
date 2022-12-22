package logger

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/tkw1536/goprogram/lib/collection"
)

// Logger is responsible for logging backups and snapshots
type Logger struct {
	component.Base
	Dependencies struct {
		SQL *sql.SQL
	}
}

// For retrieves (and prunes) the ExportLog.
// Slug determines if entries for Backups (empty slug)
// or a specific Instance (non-empty slug) are returned.
func (log *Logger) For(ctx context.Context, slug string) (exports []models.Export, err error) {
	exports, err = log.Log(ctx)
	if err != nil {
		return nil, err
	}

	return collection.Filter(exports, func(s models.Export) bool {
		return s.Slug == slug
	}), nil
}

// Log retrieves (and prunes) all entries in the snapshot log.
func (log *Logger) Log(ctx context.Context) ([]models.Export, error) {
	// query the table!
	table, err := log.Dependencies.SQL.QueryTable(ctx, false, models.ExportTable)
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
		_, err := log.Still.Environment.Stat(s.Path)
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

// AddToExportLog adds the provided export to the log.
func (log *Logger) Add(ctx context.Context, export models.Export) error {
	// find the table
	table, err := log.Dependencies.SQL.QueryTable(ctx, false, models.ExportTable)
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

// Fetch writes the SnapshotLog into the given observation
func (logger *Logger) Fetch(ctx context.Context, flags component.FetcherFlags, target *status.Distillery) (err error) {
	target.Backups, err = logger.For(ctx, "")
	return
}
