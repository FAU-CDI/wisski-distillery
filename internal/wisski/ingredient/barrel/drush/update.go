package drush

import (
	"context"
	"io"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/mstore"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/stream"
)

var errBlindUpdateFailed = exit.Error{
	Message:  "failed to run blind update script for instance %q",
	ExitCode: exit.ExitGeneric,
}

// Update performs a blind drush update
func (drush *Drush) Update(ctx context.Context, progress io.Writer) error {
	err := drush.Dependencies.Barrel.Shell(ctx, stream.NonInteractive(progress), "/runtime/blind_update.sh")
	if err != nil {
		return errBlindUpdateFailed.WithMessageF(drush.Slug).Wrap(err)
	}

	return drush.setLastUpdate(ctx)
}

const lastUpdate = mstore.For[int64]("lastUpdate")

func (drush *Drush) LastUpdate(ctx context.Context) (t time.Time, err error) {
	epoch, err := lastUpdate.Get(ctx, drush.Dependencies.MStore)
	if err == meta.ErrMetadatumNotSet {
		return t, nil
	}
	if err != nil {
		return t, err
	}

	// and turn it into time!
	return time.Unix(epoch, 0), nil
}

func (drush *Drush) setLastUpdate(ctx context.Context) error {
	return lastUpdate.Set(ctx, drush.Dependencies.MStore, time.Now().Unix())
}

type LastUpdateFetcher struct {
	ingredient.Base
	Dependencies struct {
		Drush *Drush
	}
}

var (
	_ ingredient.WissKIFetcher = (*LastUpdateFetcher)(nil)
)

func (lbr *LastUpdateFetcher) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	info.LastUpdate, err = lbr.Dependencies.Drush.LastUpdate(flags.Context)
	return
}
