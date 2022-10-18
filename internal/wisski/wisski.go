// Package wisski provides WissKI
package wisski

import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter/logger"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/triplestore"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

// WissKI represents a single WissKI Instance
type WissKI struct {
	models.Instance // whatever is stored inside the underlying instance

	// Drupal credentials - not stored in the database
	DrupalUsername string
	DrupalPassword string

	// references to components!
	Core component.Still
	Meta *meta.Meta
	TS   *triplestore.Triplestore
	SQL  *sql.SQL

	ExporterLog *logger.Logger
}

// Save saves this instance in the bookkeeping table
func (wisski *WissKI) Save() error {
	db, err := wisski.SQL.QueryTable(false, models.InstanceTable)
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
	db, err := wisski.SQL.QueryTable(false, models.InstanceTable)
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
