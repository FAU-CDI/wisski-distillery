//spellchecker:words bookkeeping
package bookkeeping

//spellchecker:words context github wisski distillery internal ingredient
import (
	"context"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
)

// Bookkeeping provides instance bookkeeping.
type Bookkeeping struct {
	ingredient.Base
}

// Save saves this instance in the bookkeeping table.
func (bk *Bookkeeping) Save(ctx context.Context) error {
	liquid := ingredient.GetLiquid(bk)
	sdb, err := ingredient.GetLiquid(bk).SQL.QueryTable(ctx, liquid.InstanceTable)
	if err != nil {
		return fmt.Errorf("failed to get bookkeeping data: %w", err)
	}

	// it has never been created => we need to create it in the database
	if liquid.Created.IsZero() {
		if err := sdb.Create(&liquid.Instance).Error; err != nil {
			return fmt.Errorf("failed to create bookkeeping data: %w", err)
		}
		return nil
	}

	// Update based on the primary key!
	if err := sdb.Select("*").Save(&liquid.Instance).Error; err != nil {
		return fmt.Errorf("failed to update bookkeeping data: %w", err)
	}
	return nil
}

// Delete deletes this instance from the bookkeeping table.
func (bk *Bookkeeping) Delete(ctx context.Context) error {
	liquid := ingredient.GetLiquid(bk)
	sdb, err := liquid.SQL.QueryTable(ctx, liquid.InstanceTable)
	if err != nil {
		return err
	}

	// doesn't exist => nothing to delete
	if liquid.Created.IsZero() {
		return nil
	}

	// delete it directly
	return sdb.Delete(&liquid.Instance).Error
}
