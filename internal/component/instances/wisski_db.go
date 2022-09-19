package instances

import "github.com/FAU-CDI/wisski-distillery/internal/models"

// WissKI represents a single WissKI Instance
type WissKI struct {
	// Whatever is stored inside the bookkeeping database
	models.Instance

	// Credentials to Drupal
	DrupalUsername string
	DrupalPassword string

	// reference to the component!
	instances *Instances
}

// Save saves this instance in the bookkeeping table
func (wisski *WissKI) Save() error {
	db, err := wisski.instances.SQL.QueryTable(false, models.InstanceTable)
	if err != nil {
		return err
	}

	// it has never been created => we need to create it in the database
	if wisski.Instance.Created.IsZero() {
		return db.Create(&wisski.Instance).Error
	}

	// Update based on the primary key!
	return db.Where("pk = ?", wisski.Instance.Pk).Updates(&wisski.Instance).Error
}

// Delete deletes this instance from the bookkeeping table
func (wisski *WissKI) Delete() error {
	db, err := wisski.instances.SQL.QueryTable(false, models.InstanceTable)
	if err != nil {
		return err
	}

	// doesn't exist => nothing to delete
	if wisski.Instance.Created.IsZero() {
		return nil
	}

	// delete it directly
	return db.Delete(&wisski.Instance).Error
}
