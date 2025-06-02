//spellchecker:words locker
package locker

//spellchecker:words context errors time github wisski distillery internal models ingredient pkglib contextx
import (
	"context"
	"errors"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/tkw1536/pkglib/contextx"
)

// Locker provides facitilites for locking this WissKI instance.
type Locker struct {
	ingredient.Base
}

var (
	_ ingredient.WissKIFetcher = (*Locker)(nil)
)

var ErrLocked = errors.New("instance is locked for administrative operations")

// TryLock attemps to lock this WissKI and returns if it suceeded.
func (lock *Locker) TryLock(ctx context.Context) bool {
	liquid := ingredient.GetLiquid(lock)

	table, err := liquid.SQL.QueryTableLegacy(ctx, liquid.LockTable)
	if err != nil {
		return false
	}

	result := table.FirstOrCreate(&models.Lock{}, models.Lock{Slug: liquid.Slug})
	return result.Error == nil && result.RowsAffected == 1
}

// TryUnlock attempts to unlock this WissKI and reports if it succeeded.
// An Unlock is also attempted when ctx is cancelled.
func (lock *Locker) TryUnlock(ctx context.Context) bool {
	liquid := ingredient.GetLiquid(lock)

	ctx, cancel := contextx.Anyways(ctx, time.Second)
	defer cancel()

	table, err := liquid.SQL.QueryTableLegacy(ctx, liquid.LockTable)
	if err != nil {
		return false
	}
	result := table.Where("slug = ?", liquid.Slug).Delete(&models.Lock{})
	return result.Error == nil && result.RowsAffected == 1
}

// Unlock unlocks this WissKI, ignoring any error.
func (lock *Locker) Unlock(ctx context.Context) {
	lock.TryUnlock(ctx)
}
