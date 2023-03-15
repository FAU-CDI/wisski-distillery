package drush

import (
	"context"
	"fmt"
	"time"

	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/stream"
)

var errCronFailed = exit.Error{
	Message:  "failed to run cron script for instance %q: exited with code %s",
	ExitCode: exit.ExitGeneric,
}

func (drush *Drush) Cron(ctx context.Context, progress io.Writer) error {
	code := drush.Dependencies.Barrel.Shell(ctx, stream.NonInteractive(progress), "/runtime/cron.sh")()
	if code != 0 {
		// keep going, because we want to run as many crons as possible
		fmt.Fprintf(progress, "%v", errCronFailed.WithMessageF(drush.Slug, code))
	}

	return nil
}

func (drush *Drush) LastCron(ctx context.Context, server *phpx.Server) (t time.Time, err error) {
	var timestamp int64
	err = drush.Dependencies.PHP.EvalCode(ctx, server, &timestamp, `$val = \Drupal::state()->get('system.cron_last'); return $val; `)
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
