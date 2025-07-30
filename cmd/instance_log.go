package cmd

//spellchecker:words github wisski distillery internal cobra pkglib errorsx exit
import (
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/exit"
)

func NewInstanceLogCommand() *cobra.Command {
	impl := new(instanceLog)

	cmd := &cobra.Command{
		Use:     "instance_log SLUG",
		Short:   "follows logs for a given instance",
		Args:    cobra.ExactArgs(1),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	return cmd
}

type instanceLog struct {
	Positionals struct {
		Slug string
	}
}

func (i *instanceLog) ParseArgs(cmd *cobra.Command, args []string) error {
	i.Positionals.Slug = args[0]
	return nil
}

var (
	errInstanceLogWissKI = exit.NewErrorWithCode("unable to get WissKI", cli.ExitGeneric)
	errInstanceLogStack  = exit.NewErrorWithCode("unable to get stack", cli.ExitGeneric)
	errInstanceLogAttach = exit.NewErrorWithCode("unable to attach to stack", cli.ExitGeneric)
)

func (i *instanceLog) Exec(cmd *cobra.Command, args []string) (e error) {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errInstanceLogWissKI, err)
	}

	instance, err := dis.Instances().WissKI(cmd.Context(), i.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: %w", errInstanceLogWissKI, err)
	}

	stack, err := instance.Barrel().OpenStack()
	if err != nil {
		return fmt.Errorf("%w: %w", errInstanceLogStack, err)
	}
	defer errorsx.Close(stack, &e, "stack")

	if err := stack.Attach(cmd.Context(), streamFromCommand(cmd), false); err != nil {
		return fmt.Errorf("%w: %w", errInstanceLogAttach, err)
	}
	return nil
}
