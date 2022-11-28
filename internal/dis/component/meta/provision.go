package meta

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

// Provision provisions new meta storage for this instance.
// NOTE(twiesing): This is a no-op, because we implement Purge.
func (meta *Meta) Provision(ctx context.Context, instance models.Instance, domain string) error {
	return nil
}

// Purge purges the storage for the given instance.
func (meta *Meta) Purge(ctx context.Context, instance models.Instance, domain string) error {
	return meta.Storage(instance.Slug).Purge(ctx)
}
