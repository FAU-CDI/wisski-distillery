package cmd

//spellchecker:words github wisski distillery internal goprogram exit pkglib status
import (
	"fmt"
	"io"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/status"
)

func NewCronCommand() *cobra.Command {
	impl := new(cron)

	cmd := &cobra.Command{
		Use:     "cron",
		Short:   "runs the cron script for several instances",
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.IntVar(&impl.Parallel, "parallel", 1, "run on (at most) this many instances in parallel. 0 for no limit.")

	return cmd
}

type cron struct {
	Parallel    int
	Positionals struct {
		Slug []string
	}
}

func (cr *cron) ParseArgs(cmd *cobra.Command, args []string) error {
	cr.Positionals.Slug = args
	return nil
}

func (*cron) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "cron",
		Description: "runs the cron script for several instances",
	}
}

var errCronFailed = exit.NewErrorWithCode("failed to run cron", exit.ExitGeneric)

func (cr *cron) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errCronFailed, err)
	}

	// find all the instances!
	wissKIs, err := dis.Instances().Load(cmd.Context(), cr.Positionals.Slug...)
	if err != nil {
		return fmt.Errorf("%w: failed to load instances: %w", errCronFailed, err)
	}

	// and do the actual blind_update!
	if err := status.WriterGroup(cmd.ErrOrStderr(), cr.Parallel, func(instance *wisski.WissKI, writer io.Writer) error {
		return instance.Drush().Cron(cmd.Context(), writer)
	}, wissKIs, status.SmartMessage(func(item *wisski.WissKI) string {
		return fmt.Sprintf("cron %q", item.Slug)
	})); err != nil {
		return fmt.Errorf("%w: %w", errCronFailed, err)
	}
	return nil
}
