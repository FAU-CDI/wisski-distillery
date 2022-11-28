package locker

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
)

// Locked checks if this WissKI is currently locked.
// If an error occurs, the instance is considered not locked.
func (lock *Locker) Locked(ctx context.Context) (locked bool) {
	table, err := lock.SQL.QueryTable(ctx, true, models.LockTable)
	if err != nil {
		return false
	}

	// check if this instance is locked
	table.Select("count(*) > 0").Where("slug = ?", lock.Slug).Find(&locked)
	return
}

func (locker *Locker) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	info.Locked = locker.Locked(flags.Context)
	return
}
