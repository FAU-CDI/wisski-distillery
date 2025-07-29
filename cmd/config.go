package cmd

//spellchecker:words github wisski distillery internal goprogram exit
import (
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
)

func NewConfigCommand() *cobra.Command {
	impl := new(cfg)

	cmd := &cobra.Command{
		Use:     "config",
		Short:   "prints information about configuration",
		Args:    cobra.NoArgs,
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.Human, "human", false, "Print configuration in human-readable format")

	return cmd
}

type cfg struct {
	Human bool
}

func (c *cfg) ParseArgs(cmd *cobra.Command, args []string) error {
	return nil
}

var errMarshalConfig = exit.NewErrorWithCode("unable to marshal config", exit.ExitGeneric)

func (c *cfg) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errMarshalConfig, err)
	}

	if c.Human {
		human := dis.Config.MarshalSensitive()
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), human)
		return nil
	}
	if err := dis.Config.Marshal(cmd.OutOrStdout()); err != nil {
		return fmt.Errorf("%w: %w", errMarshalConfig, err)
	}
	return nil
}
