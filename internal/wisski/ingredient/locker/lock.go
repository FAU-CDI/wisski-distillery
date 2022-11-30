package locker

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/tkw1536/goprogram/exit"
)

// Locker provides facitilites for locking this WissKI instance
type Locker struct {
	ingredient.Base
}

var (
	_ = (ingredient.WissKIFetcher)((*Locker)(nil))
)

var Locked = exit.Error{
	Message:  "WissKI Instance is locked for administrative operations",
	ExitCode: exit.ExitGeneric,
}

// TryLock attemps to lock this WissKI and returns if it suceeded
func (lock *Locker) TryLock(ctx context.Context) bool {
	table, err := lock.Malt.SQL.QueryTable(ctx, true, models.LockTable)
	if err != nil {
		return false
	}

	result := table.FirstOrCreate(&models.Lock{}, models.Lock{Slug: lock.Slug})
	return result.Error == nil && result.RowsAffected == 1
}

// TryUnlock attempts to unlock this WissKI and reports if it succeeded.
// An unlock can only
func (lock *Locker) TryUnlock(ctx context.Context) bool {
	table, err := lock.Malt.SQL.QueryTable(ctx, true, models.LockTable)
	if err != nil {
		return false
	}
	result := table.Where("slug = ?", lock.Slug).Delete(&models.Lock{})
	return result.Error == nil && result.RowsAffected == 1
}

// Unlock unlocks this WissKI, ignoring any error.
func (lock *Locker) Unlock(ctx context.Context) {
	lock.TryUnlock(ctx)
}
