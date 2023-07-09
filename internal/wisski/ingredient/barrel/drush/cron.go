package drush

import (
	"context"
	"fmt"
	"time"

	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/tkw1536/goprogram/exit"
)

var errCronFailed = exit.Error{
	Message:  "failed to run cron script for instance %q: exited with code %d",
	ExitCode: exit.ExitGeneric,
}

func (drush *Drush) Cron(ctx context.Context, progress io.Writer) error {
	err := drush.Exec(ctx, progress, "core-cron")
	if err != nil {
		code := err.(barrel.ExitError).Code
		// keep going, because we want to run as many crons as possible
		fmt.Fprintf(progress, "%v", errCronFailed.WithMessageF(drush.Slug, code))
	}

	return nil
}

func (drush *Drush) LastCron(ctx context.Context, server *phpx.Server) (t time.Time, err error) {
	var timestamp int64
	err = drush.Dependencies.PHP.EvalCode(ctx, server, &timestamp, `return \Drupal::state()->get('system.cron_last');`)
	if err != nil {
		return
	}
	return time.Unix(timestamp, 0), nil
}

type LastCronFetcher struct {
	ingredient.Base
	Dependencies struct {
		Drush *Drush
	}
}

var (
	_ ingredient.WissKIFetcher = (*LastCronFetcher)(nil)
)

func (lbr *LastCronFetcher) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	if flags.Quick {
		return
	}

	info.LastRebuild, _ = lbr.Dependencies.Drush.LastCron(flags.Context, flags.Server)
	return
}
