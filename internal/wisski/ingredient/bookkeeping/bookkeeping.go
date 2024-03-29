package bookkeeping

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
	sdb, err := bk.Malt.SQL.QueryTable(ctx, bk.Malt.InstanceTable)
	if err != nil {
		return err
	}

	// it has never been created => we need to create it in the database
	if bk.Instance.Created.IsZero() {
		return sdb.Create(&bk.Instance).Error
	}

	// Update based on the primary key!
	return sdb.Select("*").Save(&bk.Instance).Error
}

// Delete deletes this instance from the bookkeeping table
func (bk *Bookkeeping) Delete(ctx context.Context) error {
	sdb, err := bk.Malt.SQL.QueryTable(ctx, bk.Malt.InstanceTable)
	if err != nil {
		return err
	}

	// doesn't exist => nothing to delete
	if bk.Instance.Created.IsZero() {
		return nil
	}

	// delete it directly
	return sdb.Delete(&bk.Instance).Error
}
