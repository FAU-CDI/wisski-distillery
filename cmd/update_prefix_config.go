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
var UpdatePrefixConfig wisski_distillery.Command = updateprefixconfig{}

type updateprefixconfig struct {
	Parallel int `default:"1" description:"run on (at most) this many instances in parallel. 0 for no limit." long:"parallel" short:"p"`
}

func (updateprefixconfig) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "update_prefix_config",
		Description: "updates the prefix configuration",
	}
}

var errPrefixUpdateFailed = exit.NewErrorWithCode("failed to update prefix configuration", exit.ExitGeneric)

func (upc updateprefixconfig) Run(context wisski_distillery.Context) (err error) {
	dis := context.Environment

	wissKIs, err := dis.Instances().All(context.Context)
	if err != nil {
		return fmt.Errorf("%w: failed to get all instances: %w", errPrefixUpdateFailed, err)
	}

	if err := status.WriterGroup(context.Stderr, upc.Parallel, func(instance *wisski.WissKI, writer io.Writer) error {
		if _, err := io.WriteString(writer, "reading prefixes"); err != nil {
			return fmt.Errorf("failed to log progress: %w", err)
		}
		return instance.Prefixes().Update(context.Context)
	}, wissKIs, status.SmartMessage(func(item *wisski.WissKI) string {
		return fmt.Sprintf("update_prefix %q", item.Slug)
	})); err != nil {
		return fmt.Errorf("%w: failed to update prefixes: %w", errPrefixUpdateFailed, err)
	}
	return nil
}
