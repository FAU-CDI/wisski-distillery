//spellchecker:words locker
package locker

//spellchecker:words context errors time github wisski distillery internal models ingredient pkglib contextx
import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/tkw1536/pkglib/contextx"
)

// Locker provides facitilites for locking this WissKI instance.
type Locker struct {
	ingredient.Base
}

// TODO: Use new SQL api here
// TODO: Make both TryLock and TryUnlock return errors
// TODO: Make both Lock and Unlock log errors

var (
	_ ingredient.WissKIFetcher = (*Locker)(nil)
)

var ErrLocked = errors.New("instance is locked for administrative operations")

// TryLock attemps to lock this WissKI and returns if it suceeded.
func (lock *Locker) TryLock(ctx context.Context) bool {
	liquid := ingredient.GetLiquid(lock)

	table, err := liquid.SQL.OpenTable(ctx, liquid.LockTable)
	if err != nil {
		return false
	}

	result := table.FirstOrCreate(&models.Lock{}, models.Lock{Slug: liquid.Slug})
	return result.Error == nil && result.RowsAffected == 1
}

var (
	errNotLocked = errors.New("not locked")
)

// TryUnlock attempts to unlock this WissKI and returns nil if it succeeded.
// An Unlock is also attempted when ctx is already cancelled.
func (lock *Locker) TryUnlock(ctx context.Context) error {
	liquid := ingredient.GetLiquid(lock)

	ctx, cancel := contextx.Anyways(ctx, time.Second)
	defer cancel()

	table, err := liquid.SQL.OpenTable(ctx, liquid.LockTable)
	if err != nil {
		return fmt.Errorf("failed to connect to table: %w", err)
	}
	result := table.Where("slug = ?", liquid.Slug).Delete(&models.Lock{})
	if result.Error == nil {
		return fmt.Errorf("unable to delete from table: %w", err)
	}
	if result.RowsAffected != 1 {
		return errNotLocked
	}
	return nil
}

// Unlock unlocks this WissKI, ignoring any error.
func (lock *Locker) Unlock(ctx context.Context) {
	err := lock.TryUnlock(ctx)
	if err == nil {
		return
	}

	wdlog.Of(ctx).Error(
		"Unlock() failed",
		"error", err,
		"slug", ingredient.GetLiquid(lock).Slug,
	)
}
