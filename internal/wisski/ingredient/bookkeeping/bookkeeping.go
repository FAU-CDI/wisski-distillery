package bookkeeping

import (
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
)

// Bookkeeping provides instance bookkeeping
type Bookkeeping struct {
	ingredient.Base
}

// Save saves this instance in the bookkeeping table
func (bk *Bookkeeping) Save() error {
	sdb, err := bk.Malt.SQL.QueryTable(false, models.InstanceTable)
	if err != nil {
		return err
	}

	// it has never been created => we need to create it in the database
	if bk.Instance.Created.IsZero() {
		return sdb.Create(&bk.Instance).Error
	}

	// Update based on the primary key!
	return sdb.Where("pk = ?", bk.Instance.Pk).Updates(&bk.Instance).Error
}

// Delete deletes this instance from the bookkeeping table
func (bk *Bookkeeping) Delete() error {
	sdb, err := bk.Malt.SQL.QueryTable(false, models.InstanceTable)
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
