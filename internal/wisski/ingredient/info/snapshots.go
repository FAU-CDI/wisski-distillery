//spellchecker:words info
package info

//spellchecker:words github wisski distillery internal status ingredient
import (
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
)

type SnapshotsFetcher struct {
	ingredient.Base
}

var (
	_ ingredient.WissKIFetcher = (*SnapshotsFetcher)(nil)
)

func (lbr *SnapshotsFetcher) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	if flags.Quick {
		return
	}

	info.Snapshots, _ = ingredient.GetLiquid(lbr).Snapshots(flags.Context)
	return
}
