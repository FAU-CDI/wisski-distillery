//spellchecker:words locker
package locker

//spellchecker:words context errors time github wisski distillery internal models ingredient pkglib contextx
import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/go-sql-driver/mysql"
	"github.com/tkw1536/pkglib/contextx"
	"gorm.io/gorm"
)

// Locker provides facitilites for locking this WissKI instance.
type Locker struct {
	ingredient.Base
}

// Timeout for an unlock operation when the context is already cancelled.
// The value of this constant may change in the future.
const UnlockAnywaysTimeout = 10 * time.Second

var (
	ErrLocked    = errors.New("instance is locked for administrative operations")
	ErrNotLocked = errors.New("instance is not locked for administrative operations")
)

// TryLock attemps to lock this WissKI and returns nil if it succeeded.
// If the instance is already locked, returns an error wrapping [ErrLocked].
func (lock *Locker) TryLock(ctx context.Context) error {
	liquid := ingredient.GetLiquid(lock)

	table, err := sql.OpenInterface[models.Lock](ctx, liquid.SQL, liquid.LockTable)
	if err != nil {
		return fmt.Errorf("failed to open interface: %w", err)
	}

	{
		err := table.Create(ctx, &models.Lock{Slug: liquid.Slug})
		if isDuplicateKeyEntryError(err) {
			return fmt.Errorf("%w: %w", ErrLocked, err)
		}
		if err != nil {
			return fmt.Errorf("failed to lock instance: %w", err)
		}
	}

	return nil
}

// Error number representing a duplicate entry.
// See [mysqldocs].
//
// mysqldocs: https://dev.mysql.com/doc/mysql-errors/5.7/en/server-error-reference.html#error_er_dup_entry
const mysqlDupErrorNumber = 1062

// isDuplicateKeyEntryError checks if the given error has a duplicated primary key.
func isDuplicateKeyEntryError(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	var mysqlErr *mysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		return false
	}

	return mysqlErr.Number == mysqlDupErrorNumber
}

// TryUnlock attempts to unlock this WissKI and returns nil if it succeeded.
// As a special case to avoid deadlocks, an unlock is also attempted when ctx is already cancelled.
// In such a case, the timeout for the unlock is [UnlockAnywaysTimeout].
func (lock *Locker) TryUnlock(ctx context.Context) error {
	ctx, cancel := contextx.Anyways(ctx, UnlockAnywaysTimeout)
	defer cancel()

	liquid := ingredient.GetLiquid(lock)
	table, err := sql.OpenInterface[models.Lock](ctx, liquid.SQL, liquid.LockTable)
	if err != nil {
		return fmt.Errorf("failed to open interface: %w", err)
	}

	count, err := table.Where("slug = ?", liquid.Slug).Delete(ctx)
	if err != nil {
		return fmt.Errorf("unable to delete from table: %w", err)
	}
	if count == 0 {
		return ErrNotLocked
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
