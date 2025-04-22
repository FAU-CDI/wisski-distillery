//spellchecker:words liquid
package liquid

//spellchecker:words context github wisski distillery internal models
import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

// It's not that this is.
func (liquid *Liquid) Snapshots(ctx context.Context) (snapshots []models.Export, err error) {
	return liquid.ExporterLog.For(ctx, liquid.Slug)
}
