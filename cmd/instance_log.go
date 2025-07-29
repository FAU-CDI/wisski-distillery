package cmd

//spellchecker:words github wisski distillery internal goprogram exit pkglib errorsx
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/exit"
)

func NewInstanceLogCommand() *cobra.Command {
	impl := new(instanceLog)

	cmd := &cobra.Command{
		Use:     "instance_log",
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

func (*instanceLog) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "instance_log",
		Description: "follows logs for a given instance",
	}
}

var (
	errInstanceLogWissKI = exit.NewErrorWithCode("unable to get WissKI", exit.ExitGeneric)
	errInstanceLogStack  = exit.NewErrorWithCode("unable to get stack", exit.ExitGeneric)
	errInstanceLogAttach = exit.NewErrorWithCode("unable to attach to stack", exit.ExitGeneric)
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
