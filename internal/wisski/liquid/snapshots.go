package liquid

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

// Snapshots returns the list of snapshots of this WissKI
// NOTE(twiesing): Not entirely sure where this should go.
// It's not that this is
func (liquid *Liquid) Snapshots(ctx context.Context) (snapshots []models.Export, err error) {
	return liquid.Malt.ExporterLog.For(ctx, liquid.Slug)
}
