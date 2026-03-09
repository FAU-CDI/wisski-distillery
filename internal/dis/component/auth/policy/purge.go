//spellchecker:words policy
package policy

//spellchecker:words context github wisski distillery internal models
import (
	"context"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

func (*Policy) ProvisionNeedsStack(instance models.Instance) bool {
	return false
}

func (*Policy) Provision(ctx context.Context, instance models.Instance, domain string, stack *component.StackWithResources) error {
	// component is purge-only
	return nil
}

func (*Policy) PurgeMayFail(instance models.Instance) bool {
	return false
}

// Purge purges every policy for the given slug form the database.
func (pol *Policy) Purge(ctx context.Context, instance models.Instance, domain string) error {
	table, err := pol.openInterface(ctx)
	if err != nil {
		return err
	}
	if _, err := table.Where(&models.Grant{Slug: instance.Slug}).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	return nil
}

// OnUserDelete is called when a user is deleted.
func (pol *Policy) OnUserDelete(ctx context.Context, user *models.User) error {
	table, err := pol.openInterface(ctx)
	if err != nil {
		return err
	}

	if _, err := table.Where(&models.Grant{User: user.User}).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	return nil
}
