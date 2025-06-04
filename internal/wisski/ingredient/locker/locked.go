//spellchecker:words locker
package locker

//spellchecker:words context github wisski distillery internal status ingredient
import (
	"context"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/liquid"
)

// Locked checks if this WissKI is currently locked.
// If an error occurs, the instance is considered not locked.
func (lock *Locker) Locked(ctx context.Context) (locked bool) {
	liquid := ingredient.GetLiquid(lock)
	locked, err := lock.locked(ctx, liquid)
	if err != nil {
		wdlog.Of(ctx).Error(
			"failed to check for locked state, returning false",
			"slug", liquid.Slug,
			"error", err,
		)
		return false
	}
	return locked
}

func (lock *Locker) locked(ctx context.Context, liquid *liquid.Liquid) (locked bool, err error) {
	table, err := sql.OpenInterface[models.Lock](ctx, liquid.SQL, liquid.LockTable)
	if err != nil {
		return false, fmt.Errorf("failed to open interface: %w", err)
	}
	res, err := table.Where("slug = ?", liquid.Slug).Count(ctx, "*")
	if err != nil {
		return false, fmt.Errorf("failed to query table: %w", err)
	}
	return res > 0, nil
}

func (locker *Locker) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	info.Locked = locker.Locked(flags.Context)
	return
}
