//spellchecker:words drush
package drush

//spellchecker:words context errors time github wisski distillery internal phpx status ingredient barrel
import (
	"context"
	"errors"
	"fmt"
	"time"

	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
)

func (drush *Drush) Cron(ctx context.Context, progress io.Writer) error {
	err := drush.Exec(ctx, progress, "core-cron")
	if err != nil {
		var ee barrel.ExitError
		if !(errors.As(err, &ee)) {
			return fmt.Errorf("drush.Exec returned unexpected error: %w", err)
		}
		code := ee.Code()

		// keep going, because we want to run as many crons as possible
		if _, err := fmt.Fprintf(progress, "failed to run cron script for instance %q: exited with code %d", ingredient.GetLiquid(drush).Slug, code); err != nil {
			return fmt.Errorf("failed to report progress: %w", err)
		}
	}

	return nil
}

func (drush *Drush) LastCron(ctx context.Context, server *phpx.Server) (t time.Time, err error) {
	var timestamp int64
	err = drush.dependencies.PHP.EvalCode(ctx, server, &timestamp, `return \Drupal::state()->get('system.cron_last');`)
	if err != nil {
		return
	}
	return time.Unix(timestamp, 0), nil
}

type LastCronFetcher struct {
	ingredient.Base
	dependencies struct {
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

	info.LastCron, _ = lbr.dependencies.Drush.LastCron(flags.Context, flags.Server)
	return
}
