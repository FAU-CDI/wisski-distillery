package cmd

//spellchecker:words github wisski distillery internal component exporter goprogram exit
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
)

func NewSnapshotCommand() *cobra.Command {
	impl := new(snapshot)

	cmd := &cobra.Command{
		Use:     "snapshot",
		Short:   "generates a snapshot archive for the provided instance",
		Args:    cobra.RangeArgs(1, 2),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.Keepalive, "keepalive", false, "keep instance running while taking a backup. might lead to inconsistent state")
	flags.BoolVar(&impl.StagingOnly, "staging-only", false, "do not package into a snapshot archive, but only create a staging directory")
	flags.StringSliceVar(&impl.Parts, "parts", nil, "parts to include in snapshots. defaults to all parts, use l to list all available parts")
	flags.BoolVar(&impl.List, "list-parts", false, "list available parts")

	return cmd
}

type snapshot struct {
	Keepalive   bool
	StagingOnly bool
	Parts       []string
	List        bool
	Positionals struct {
		Slug string
		Dest string
	}
}

func (sn *snapshot) ParseArgs(cmd *cobra.Command, args []string) error {
	sn.Positionals.Slug = args[0]
	if len(args) >= 2 {
		sn.Positionals.Dest = args[1]
	}
	return nil
}

func (*snapshot) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "snapshot",
		Description: "generates a snapshot archive for the provided instance",
	}
}

var (
	errSnapshotFailed = exit.NewErrorWithCode("failed to make a snapshot", exit.ExitGeneric)
	errSnapshotWissKI = exit.NewErrorWithCode("unable to find WissKI", exit.ExitGeneric)
)

func (sn *snapshot) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errSnapshotFailed, err)
	}

	// list available parts
	if sn.List {
		for _, part := range dis.Exporter().Parts() {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), part)
		}
		return nil
	}

	// find the instance!
	instance, err := dis.Instances().WissKI(cmd.Context(), sn.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: %w", errSnapshotWissKI, err)
	}

	// do a snapshot of it!
	err = dis.Exporter().MakeExport(cmd.Context(), cmd.ErrOrStderr(), exporter.ExportTask{
		Dest:        sn.Positionals.Dest,
		StagingOnly: sn.StagingOnly,

		SnapshotDescription: exporter.SnapshotDescription{
			Parts: sn.Parts,
		},
		Instance: instance,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errSnapshotFailed, err)
	}
	return nil
}
