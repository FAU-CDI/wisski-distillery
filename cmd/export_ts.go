package cmd

//spellchecker:words github wisski distillery internal cobra
import (
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
)

func NewExportTSCommand() *cobra.Command {
	impl := new(exportTS)

	cmd := &cobra.Command{
		Use:     "export_ts SLUG",
		Short:   "export the triplestore for a specific instance",
		Args:    cobra.ExactArgs(1),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	return cmd
}

type exportTS struct {
	Positionals struct {
		Slug string
	}
}

func (ets *exportTS) ParseArgs(cmd *cobra.Command, args []string) error {
	if len(args) >= 1 {
		ets.Positionals.Slug = args[0]
	}
	return nil
}

func (ets *exportTS) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("failed to get WissKI: %w", err)
	}

	instance, err := dis.Instances().WissKI(cmd.Context(), ets.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("failed to get WissKI: %w", err)
	}

	if err := instance.BoundTriplestore().SnapshotDB(cmd.Context(), cmd.OutOrStdout()); err != nil {
		return fmt.Errorf("failed to snapshot triplestore: %w", err)
	}
	return nil
}
