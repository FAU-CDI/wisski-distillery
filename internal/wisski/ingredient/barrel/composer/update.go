package composer

import (
	"context"
	"io"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/mstore"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
)

var errBlindUpdateFailed = exit.Error{
	Message:  "failed to run blind update script for instance %q",
	ExitCode: exit.ExitGeneric,
}

// Update performs a blind drush update
func (composer *Composer) Update(ctx context.Context, progress io.Writer) (err error) {
	defer errBlindUpdateFailed.WithMessageF(composer.Slug).DeferWrap(&err)

	if err := composer.FixPermission(ctx, progress); err != nil {
		return err
	}

	logging.LogMessage(progress, "Updating Packages")
	{
		err := composer.Exec(ctx, progress, "update")
		if err != nil {
			return err
		}
	}

	logging.LogMessage(progress, "Installing database updates")
	{
		err := composer.dependencies.Drush.Exec(ctx, progress, "-y", "updatedb")
		if err != nil {
			return err
		}
	}

	logging.LogMessage(progress, "Updating WissKI Packages")
	{
		err := composer.Exec(ctx, progress, "update")
		if err != nil {
			return err
		}
	}

	return composer.setLastUpdate(ctx)
}

const lastUpdate = mstore.For[int64]("lastUpdate")

func (drush *Composer) LastUpdate(ctx context.Context) (t time.Time, err error) {
	epoch, err := lastUpdate.Get(ctx, drush.dependencies.MStore)
	if err == meta.ErrMetadatumNotSet {
		return t, nil
	}
	if err != nil {
		return t, err
	}

	// and turn it into time!
	return time.Unix(epoch, 0), nil
}

func (drush *Composer) setLastUpdate(ctx context.Context) error {
	return lastUpdate.Set(ctx, drush.dependencies.MStore, time.Now().Unix())
}

type LastUpdateFetcher struct {
	ingredient.Base
	dependencies struct {
		Composer *Composer
	}
}

var (
	_ ingredient.WissKIFetcher = (*LastUpdateFetcher)(nil)
)

func (lbr *LastUpdateFetcher) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	info.LastUpdate, err = lbr.dependencies.Composer.LastUpdate(flags.Context)
	return
}
