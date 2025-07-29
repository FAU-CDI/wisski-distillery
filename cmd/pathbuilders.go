package cmd

//spellchecker:words github wisski distillery internal cobra pkglib exit
import (
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
)

func NewPathbuildersCommand() *cobra.Command {
	impl := new(pathbuilders)

	cmd := &cobra.Command{
		Use:     "pathbuilders SLUG [NAME]",
		Short:   "list pathbuilders of a specific instance",
		Args:    cobra.RangeArgs(1, 2),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	return cmd
}

type pathbuilders struct {
	Positionals struct {
		Slug string
		Name string
	}
}

func (pb *pathbuilders) ParseArgs(cmd *cobra.Command, args []string) error {
	pb.Positionals.Slug = args[0]
	if len(args) >= 2 {
		pb.Positionals.Name = args[1]
	}
	return nil
}

var (
	errPathbuildersExport  = exit.NewErrorWithCode("unable to export pathbuilder", exit.ExitGeneric)
	errPathbuildersNoExist = exit.NewErrorWithCode("pathbuilder does not exist", exit.ExitGeneric)
	errPathbuildersWissKI  = exit.NewErrorWithCode("unable to find WissKI", exit.ExitGeneric)
)

func (pb *pathbuilders) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errPathbuildersWissKI, err)
	}

	// get the wisski
	instance, err := dis.Instances().WissKI(cmd.Context(), pb.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: %w", errPathbuildersWissKI, err)
	}

	// get all of the pathbuilders
	if pb.Positionals.Name == "" {
		names, err := instance.Pathbuilder().All(cmd.Context(), nil)
		if err != nil {
			return fmt.Errorf("%w: %w", errPathbuildersExport, err)
		}
		for _, name := range names {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), name)
		}
		return nil
	}

	// get all the pathbuilders
	xml, err := instance.Pathbuilder().Get(cmd.Context(), nil, pb.Positionals.Name)
	if xml == "" {
		return fmt.Errorf("%q: %w", pb.Positionals.Name, errPathbuildersNoExist)
	}
	if err != nil {
		return fmt.Errorf("%w: %w", errPathbuildersExport, err)
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s", xml)

	return nil
}
