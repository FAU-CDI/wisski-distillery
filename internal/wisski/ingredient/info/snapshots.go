package info

import (
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
)

type SnapshotsFetcher struct {
	ingredient.Base

	Info *Info
}

func (lbr *SnapshotsFetcher) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	if flags.Quick {
		return
	}

	info.Snapshots, _ = lbr.Snapshots()
	return
}
