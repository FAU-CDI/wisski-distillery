//spellchecker:words liquid
package liquid

//spellchecker:words context github wisski distillery internal models
import (
	"context"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

func (liquid *Liquid) Snapshots(ctx context.Context) (snapshots []models.Export, err error) {
	snapshots, err = liquid.ExporterLog.For(ctx, liquid.Slug)
	if err != nil {
		return nil, fmt.Errorf("failed to list snapshots: %w", err)
	}
	return snapshots, nil
}
