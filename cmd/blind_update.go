package cmd

//spellchecker:words github wisski distillery internal goprogram exit pkglib collection status
import (
	"fmt"
	"io"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/collection"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/status"
)

func NewBlindUpdateCommand() *cobra.Command {
	impl := new(blindUpdate)

	cmd := &cobra.Command{
		Use:     "blind_update",
		Short:   "runs the blind update in the provided instances",
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.IntVar(&impl.Parallel, "parallel", 1, "run on (at most) this many instances in parallel. 0 for no limit")
	flags.BoolVar(&impl.Force, "force", false, "force running blind-update even if 'AutoBlindUpdate' is set to false")

	return cmd
}

type blindUpdate struct {
	Parallel    int
	Force       bool
	Positionals struct {
		Slug []string
	}
}

func (bu *blindUpdate) ParseArgs(cmd *cobra.Command, args []string) error {
	bu.Positionals.Slug = args
	return nil
}

func (*blindUpdate) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "blind_update",
		Description: "runs the blind update in the provided instances",
	}
}

var errBlindUpdateFailed = exit.NewErrorWithCode("failed to run blind update", exit.ExitGeneric)

func (bu *blindUpdate) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errBlindUpdateFailed, err)
	}

	// find all the instances!
	wissKIs, err := dis.Instances().Load(cmd.Context(), bu.Positionals.Slug...)
	if err != nil {
		return fmt.Errorf("%w: %w", errBlindUpdateFailed, err)
	}
	if !bu.Force {
		wissKIs = collection.KeepFunc(wissKIs, func(instance *wisski.WissKI) bool {
			return bool(instance.AutoBlindUpdateEnabled)
		})
	}

	// and do the actual blind_update!
	if err := status.WriterGroup(cmd.ErrOrStderr(), bu.Parallel, func(instance *wisski.WissKI, writer io.Writer) error {
		return instance.Composer().Update(cmd.Context(), writer)
	}, wissKIs, status.SmartMessage(func(item *wisski.WissKI) string {
		return fmt.Sprintf("blind_update %q", item.Slug)
	})); err != nil {
		return fmt.Errorf("%w: failed to blind_update: %w", errBlindUpdateFailed, err)
	}
	return nil
}
