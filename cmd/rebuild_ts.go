package cmd

//spellchecker:words github wisski distillery internal
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
)

func NewRebuildTSCommand() *cobra.Command {
	impl := new(rebuildTS)

	cmd := &cobra.Command{
		Use:     "rebuild_ts SLUG",
		Short:   "rebuild the triplestore for a specific instance",
		Args:    cobra.ExactArgs(1),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.AllowEmptyRepository, "allow-empty", false, "don't abort if repository is empty")

	return cmd
}

type rebuildTS struct {
	AllowEmptyRepository bool
	Positionals          struct {
		Slug string
	}
}

func (rts *rebuildTS) ParseArgs(cmd *cobra.Command, args []string) error {
	if len(args) >= 1 {
		rts.Positionals.Slug = args[0]
	}
	return nil
}

func (*rebuildTS) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "rebuild_ts",
		Description: "rebuild the triplestore for a specific instance",
	}
}

func (rts *rebuildTS) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("failed to get WissKI: %w", err)
	}

	instance, err := dis.Instances().WissKI(cmd.Context(), rts.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("failed to get WissKI: %w", err)
	}

	_, err = instance.TRB().RebuildTriplestore(cmd.Context(), cmd.OutOrStdout(), rts.AllowEmptyRepository)
	if err != nil {
		return fmt.Errorf("failed to rebuild triplestore: %w", err)
	}
	return nil
}
