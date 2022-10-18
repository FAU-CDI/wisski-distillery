package locker

import "github.com/FAU-CDI/wisski-distillery/internal/models"

// Locked checks if this WissKI is currently locked.
func (lock *Locker) Locked() (locked bool) {
	table, err := lock.SQL.QueryTable(true, models.LockTable)
	if err != nil {
		return false
	}

	// check if this instance is locked
	table.Select("count(*) > 0").Where("slug = ?", lock.Slug).Find(&locked)
	return
}
