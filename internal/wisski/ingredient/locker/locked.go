//spellchecker:words locker
package locker

//spellchecker:words context github wisski distillery internal status ingredient
import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
)

// Locked checks if this WissKI is currently locked.
// If an error occurs, the instance is considered not locked.
func (lock *Locker) Locked(ctx context.Context) (locked bool) {
	liquid := ingredient.GetLiquid(lock)

	table, err := liquid.SQL.QueryTableLegacy(ctx, liquid.LockTable)
	if err != nil {
		return false
	}

	// check if this instance is locked
	table.Select("count(*) > 0").Where("slug = ?", liquid.Slug).Find(&locked)
	return
}

func (locker *Locker) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	info.Locked = locker.Locked(flags.Context)
	return
}
