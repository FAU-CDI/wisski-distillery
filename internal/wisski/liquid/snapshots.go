package liquid

import "github.com/FAU-CDI/wisski-distillery/internal/models"

// Snapshots returns the list of snapshots of this WissKI
// NOTE(twiesing): Not entirely sure where this should go.
// It's not that this is
func (liquid *Liquid) Snapshots() (snapshots []models.Export, err error) {
	return liquid.Malt.ExporterLog.For(liquid.Slug)
}
