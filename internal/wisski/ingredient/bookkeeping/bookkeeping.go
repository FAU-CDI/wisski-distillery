//spellchecker:words bookkeeping
package bookkeeping

//spellchecker:words context github wisski distillery internal ingredient
import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
)

// Bookkeeping provides instance bookkeeping
type Bookkeeping struct {
	ingredient.Base
}

// Save saves this instance in the bookkeeping table
func (bk *Bookkeeping) Save(ctx context.Context) error {
	liquid := ingredient.GetLiquid(bk)
	sdb, err := ingredient.GetLiquid(bk).Malt.SQL.QueryTable(ctx, liquid.Malt.InstanceTable)
	if err != nil {
		return err
	}

	// it has never been created => we need to create it in the database
	if liquid.Instance.Created.IsZero() {
		return sdb.Create(&liquid.Instance).Error
	}

	// Update based on the primary key!
	return sdb.Select("*").Save(&liquid.Instance).Error
}

// Delete deletes this instance from the bookkeeping table
func (bk *Bookkeeping) Delete(ctx context.Context) error {
	liquid := ingredient.GetLiquid(bk)
	sdb, err := liquid.SQL.QueryTable(ctx, liquid.InstanceTable)
	if err != nil {
		return err
	}

	// doesn't exist => nothing to delete
	if liquid.Instance.Created.IsZero() {
		return nil
	}

	// delete it directly
	return sdb.Delete(&liquid.Instance).Error
}
