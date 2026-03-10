//spellchecker:words meta
package meta

//spellchecker:words context github wisski distillery internal models
import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

func (meta *Meta) ProvisionNeedsStack(instance models.Instance) bool {
	return false
}

// Provision provisions new meta storage for this instance.
// NOTE(twiesing): This is a no-op, because we implement Purge.
func (meta *Meta) Provision(ctx context.Context, progress io.Writer, instance models.Instance, domain string, stack *component.StackWithResources) error {
	return nil
}

func (*Meta) PurgeMayFail(instance models.Instance) bool {
	return false
}

// Purge purges the storage for the given instance.
func (meta *Meta) Purge(ctx context.Context, progress io.Writer, instance models.Instance, domain string) error {
	return meta.Storage(instance.Slug).Purge(ctx)
}
