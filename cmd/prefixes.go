package cmd

//spellchecker:words github wisski distillery internal cobra pkglib exit
import (
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
)

func NewPrefixesCommand() *cobra.Command {
	impl := new(prefixes)

	cmd := &cobra.Command{
		Use:     "prefixes SLUG",
		Short:   "list all prefixes for a specific instance",
		Args:    cobra.ExactArgs(1),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	return cmd
}

type prefixes struct {
	Positionals struct {
		Slug string
	}
}

func (p *prefixes) ParseArgs(cmd *cobra.Command, args []string) error {
	p.Positionals.Slug = args[0]
	return nil
}

var (
	errPrefixesGeneric = exit.NewErrorWithCode("unable to load prefixes", cli.ExitGeneric)
	errPrefixesWissKI  = exit.NewErrorWithCode("unable to find WissKI", cli.ExitGeneric)
)

func (p *prefixes) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errPrefixesWissKI, err)
	}

	instance, err := dis.Instances().WissKI(cmd.Context(), p.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: %w", errPrefixesWissKI, err)
	}

	prefixes, err := instance.Prefixes().All(cmd.Context(), nil)
	if err != nil {
		return fmt.Errorf("%w: %w", errPrefixesGeneric, err)
	}

	for _, prefix := range prefixes {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), prefix)
	}

	return nil
}
