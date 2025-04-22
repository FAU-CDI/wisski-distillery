//spellchecker:words policy
package policy

//spellchecker:words context github wisski distillery internal models
import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

func (*Policy) Provision(ctx context.Context, instance models.Instance, domain string) error {
	// component is purge-only
	return nil
}

// Purge purges every policy for the given slug form the database.
func (pol *Policy) Purge(ctx context.Context, instance models.Instance, domain string) error {
	table, err := pol.table(ctx)
	if err != nil {
		return err
	}
	return table.Delete(&models.Grant{}, &models.Grant{Slug: instance.Slug}).Error
}

// OnUserDelete is called when a user is deleted.
func (pol *Policy) OnUserDelete(ctx context.Context, user *models.User) error {
	table, err := pol.table(ctx)
	if err != nil {
		return err
	}
	return table.Delete(&models.Grant{}, &models.Grant{User: user.User}).Error
}
