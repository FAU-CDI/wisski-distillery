//spellchecker:words logger
package logger

//spellchecker:words context errors reflect github wisski distillery internal component models status pkglib collection
import (
	"context"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/tkw1536/pkglib/collection"
)

// Logger is responsible for logging backups and snapshots.
type Logger struct {
	component.Base
	dependencies struct {
		SQL *sql.SQL
	}
}

var (
	_ component.Table = (*Logger)(nil)
)

func (*Logger) TableInfo() component.TableInfo {
	return component.TableInfo{
		Model: models.Export{},
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

	return collection.KeepFunc(exports, func(s models.Export) bool {
		return s.Slug == slug
	}), nil
}

// Log retrieves and cleans up all entries in the snapshot log.
func (log *Logger) Log(ctx context.Context) ([]models.Export, error) {
	table, err := sql.OpenInterface[models.Export](ctx, log.dependencies.SQL, log)
	if err != nil {
		return nil, fmt.Errorf("failed to open interface: %w", err)
	}

	// find all the exports
	exports, err := table.Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve existing exports: %w", err)
	}

	// partition out the exports that no longer exist!
	parts := collection.Partition(exports, func(s models.Export) bool {
		exists, err := s.Exists()
		if err == nil {
			return exists
		}

		wdlog.Of(ctx).Error(
			"unable to check if export exists, skipping pruning",
			"error", err,
			"pk", s.Pk,
		)
		return true
	})

	// delete the parts which no longer exist
	if len(parts[false]) > 0 {
		pks := collection.MapSlice(parts[false], func(s models.Export) uint { return s.Pk })
		if _, err := table.Where("pk in ?", pks).Delete(ctx); err != nil {
			return nil, fmt.Errorf("failed to remove old export entries: %w", err)
		}
	}

	return parts[true], nil
}

// AddToExportLog adds the provided export to the log.
func (log *Logger) Add(ctx context.Context, export models.Export) error {
	table, err := sql.OpenInterface[models.Export](ctx, log.dependencies.SQL, log)
	if err != nil {
		return fmt.Errorf("failed to open interface: %w", err)
	}

	if err := table.Create(ctx, &export); err != nil {
		return fmt.Errorf("failed to add log: %w", err)
	}
	return nil
}

// Fetch writes the SnapshotLog into the given observation.
func (logger *Logger) Fetch(ctx context.Context, flags component.FetcherFlags, target *status.Distillery) (err error) {
	target.Backups, err = logger.For(ctx, "")
	return
}
