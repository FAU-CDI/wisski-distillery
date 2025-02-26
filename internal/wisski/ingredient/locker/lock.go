//spellchecker:words locker
package locker

//spellchecker:words context time github wisski distillery internal models ingredient goprogram exit pkglib contextx
import (
	"context"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/contextx"
)

// Locker provides facitilites for locking this WissKI instance
type Locker struct {
	ingredient.Base
}

var (
	_ ingredient.WissKIFetcher = (*Locker)(nil)
)

var Locked = exit.Error{
	Message:  "instance is locked for administrative operations",
	ExitCode: exit.ExitGeneric,
}

// TryLock attemps to lock this WissKI and returns if it suceeded
func (lock *Locker) TryLock(ctx context.Context) bool {
	liquid := ingredient.GetLiquid(lock)

	table, err := liquid.SQL.QueryTable(ctx, liquid.LockTable)
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

	ctx, close := contextx.Anyways(ctx, time.Second)
	defer close()

	table, err := liquid.SQL.QueryTable(ctx, liquid.LockTable)
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
