package cmd

//spellchecker:words github wisski distillery internal goprogram exit pkglib status
import (
	"fmt"
	"io"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"go.tkw01536.de/goprogram/exit"
	"go.tkw01536.de/pkglib/status"
)

// Cron is the 'cron' command.
var Cron wisski_distillery.Command = cron{}

type cron struct {
	Parallel int `default:"1" description:"run on (at most) this many instances in parallel. 0 for no limit." long:"parallel" short:"p"`

	Positionals struct {
		Slug []string `description:"slug of instances to run cron in" positional-arg-name:"SLUG" required:"0"`
	} `positional-args:"true"`
}

func (cron) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "cron",
		Description: "runs the cron script for several instances",
	}
}

var errCronFailed = exit.NewErrorWithCode("failed to run cron", exit.ExitGeneric)

func (cr cron) Run(context wisski_distillery.Context) (err error) {
	// find all the instances!
	wissKIs, err := context.Environment.Instances().Load(context.Context, cr.Positionals.Slug...)
	if err != nil {
		return fmt.Errorf("%w: failed to load instances: %w", errCronFailed, err)
	}

	// and do the actual blind_update!
	if err := status.WriterGroup(context.Stderr, cr.Parallel, func(instance *wisski.WissKI, writer io.Writer) error {
		return instance.Drush().Cron(context.Context, writer)
	}, wissKIs, status.SmartMessage(func(item *wisski.WissKI) string {
		return fmt.Sprintf("cron %q", item.Slug)
	})); err != nil {
		return fmt.Errorf("%w: %w", errCronFailed, err)
	}
	return nil
}
