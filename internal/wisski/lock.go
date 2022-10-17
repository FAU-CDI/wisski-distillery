package wisski

import (
	"errors"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

var ErrLocked = errors.New("instance is locked")

// TryLock attemps to lock this WissKI
// If this is not possible, returns ErrLocked
func (wisski *WissKI) TryLock() error {
	table, err := wisski.SQL.QueryTable(true, models.LockTable)
	if err != nil {
		return ErrLocked
	}

	result := table.FirstOrCreate(&models.Lock{}, models.Lock{Slug: wisski.Slug})
	locked := result.Error == nil && result.RowsAffected == 1

	if !locked {
		return ErrLocked
	}
	return nil
}

func (wisski *WissKI) IsLocked() (locked bool) {
	table, err := wisski.SQL.QueryTable(true, models.LockTable)
	if err != nil {
		return false
	}

	// check if this instance is locked
	table.Select("count(*) > 0").Where("slug = ?", wisski.Slug).Find(&locked)
	return
}

// Unlock unlocks this WissKI instance and returns if it succeeded
func (wisski WissKI) Unlock() bool {
	table, err := wisski.SQL.QueryTable(true, models.LockTable)
	if err != nil {
		return false
	}
	result := table.Where("slug = ?", wisski.Slug).Delete(&models.Lock{})
	return result.Error == nil && result.RowsAffected == 1
}
