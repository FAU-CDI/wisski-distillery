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

func NewUpdatePrefixConfigCommand() *cobra.Command {
	impl := new(updateprefixconfig)

	cmd := &cobra.Command{
		Use:     "update_prefix_config",
		Short:   "updates the prefix configuration",
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.IntVar(&impl.Parallel, "parallel", 1, "run on (at most) this many instances in parallel. 0 for no limit.")

	return cmd
}

type updateprefixconfig struct {
	Parallel int
}

func (upc *updateprefixconfig) ParseArgs(cmd *cobra.Command, args []string) error {
	return nil
}

func (*updateprefixconfig) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "update_prefix_config",
		Description: "updates the prefix configuration",
	}
}

var errPrefixUpdateFailed = exit.NewErrorWithCode("failed to update prefix configuration", exit.ExitGeneric)

func (upc *updateprefixconfig) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("%w: failed to get all instances: %w", errPrefixUpdateFailed, err)
	}

	wissKIs, err := dis.Instances().All(cmd.Context())
	if err != nil {
		return fmt.Errorf("%w: failed to get all instances: %w", errPrefixUpdateFailed, err)
	}

	if err := status.WriterGroup(cmd.ErrOrStderr(), upc.Parallel, func(instance *wisski.WissKI, writer io.Writer) error {
		if _, err := io.WriteString(writer, "reading prefixes"); err != nil {
			return fmt.Errorf("failed to log progress: %w", err)
		}
		return instance.Prefixes().Update(cmd.Context())
	}, wissKIs, status.SmartMessage(func(item *wisski.WissKI) string {
		return fmt.Sprintf("update_prefix %q", item.Slug)
	})); err != nil {
		return fmt.Errorf("%w: failed to update prefixes: %w", errPrefixUpdateFailed, err)
	}
	return nil
}
