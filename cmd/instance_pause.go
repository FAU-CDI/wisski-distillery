package cmd

//spellchecker:words github wisski distillery internal goprogram exit pkglib errorsx
import (
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/exit"
)

func NewInstancePauseCommand() *cobra.Command {
	impl := new(instancepause)

	cmd := &cobra.Command{
		Use:     "instance_pause SLUG",
		Short:   "stops or starts a single instance",
		Args:    cobra.ExactArgs(1),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.Stop, "stop", false, "stop instance")
	flags.BoolVar(&impl.Start, "start", false, "start (or restart) instance")

	return cmd
}

type instancepause struct {
	Stop        bool
	Start       bool
	Positionals struct {
		Slug string
	}
}

func (i *instancepause) ParseArgs(cmd *cobra.Command, args []string) error {
	i.Positionals.Slug = args[0]

	if i.Stop == i.Start {
		return errStopStartExcluded
	}
	return nil
}

var errStopStartExcluded = exit.NewErrorWithCode("stop and start are mutually exclusive", exit.ExitCommandArguments)
var errInstancePauseWissKI = exit.NewErrorWithCode("unable to get WissKI", exit.ExitGeneric)
var errInstancePauseStack = exit.NewErrorWithCode("unable to get stack", exit.ExitGeneric)

func (i *instancepause) Exec(cmd *cobra.Command, args []string) (e error) {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errInstancePauseWissKI, err)
	}

	instance, err := dis.Instances().WissKI(cmd.Context(), i.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: %w", errInstancePauseWissKI, err)
	}

	stack, err := instance.Barrel().OpenStack()
	if err != nil {
		return fmt.Errorf("%w: %w", errInstancePauseStack, err)
	}
	defer errorsx.Close(stack, &e, "stack")

	if i.Stop {
		if err := stack.Down(cmd.Context(), cmd.OutOrStdout()); err != nil {
			return fmt.Errorf("failed to stop instance: %w", err)
		}
	} else {
		if err := stack.Start(cmd.Context(), cmd.OutOrStdout()); err != nil {
			return fmt.Errorf("failed to start instance: %w", err)
		}
	}
	return nil
}
