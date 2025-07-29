package cmd

//spellchecker:words github wisski distillery internal component logging goprogram exit pkglib errorsx status
import (
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/status"
)

func NewSystemPauseCommand() *cobra.Command {
	impl := new(systempause)

	cmd := &cobra.Command{
		Use:     "system_pause",
		Short:   "stops or starts the entire distillery system",
		Args:    cobra.NoArgs,
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.Stop, "stop", false, "stop all the components")
	flags.BoolVar(&impl.Start, "start", false, "start all the components")

	return cmd
}

type systempause struct {
	Stop  bool
	Start bool
}

func (sp *systempause) ParseArgs(cmd *cobra.Command, args []string) error {
	if sp.Stop == sp.Start {
		return errPauseArguments
	}
	return nil
}

var (
	errPauseGeneric   = exit.NewErrorWithCode("unable to pause or resume system", exit.ExitGeneric)
	errPauseArguments = exit.NewErrorWithCode("exactly one of `--stop` and `--start` must be provided", exit.ExitCommandArguments)
)

func (sp *systempause) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errPauseGeneric, err)
	}

	if sp.Start {
		err = sp.start(cmd, dis)
	} else {
		err = sp.stop(cmd, dis)
	}

	if err != nil {
		return fmt.Errorf("%w: %w", errPauseGeneric, err)
	}
	return nil
}

func (sp *systempause) start(cmd *cobra.Command, dis *dis.Distillery) error {
	if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Starting Components"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}

	// find all the core stacks
	if err := status.RunErrorGroup(cmd.ErrOrStderr(), status.Group[component.Installable, error]{
		PrefixString: func(item component.Installable, index int) string {
			return fmt.Sprintf("[up %q]: ", item.Name())
		},
		PrefixAlign: true,
		Handler: func(item component.Installable, index int, writer io.Writer) (e error) {
			stack, err := item.OpenStack()
			if err != nil {
				return fmt.Errorf("failed to open stack: %w", err)
			}
			defer errorsx.Close(stack, &e, "stack")
			return stack.Start(cmd.Context(), writer)
		},
	}, dis.Installable()); err != nil {
		return fmt.Errorf("failed to start components: %w", err)
	}

	if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Starting Up WissKIs"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}

	wissKIs, err := dis.Instances().All(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to get all instances: %w", err)
	}

	if err := status.RunErrorGroup(cmd.ErrOrStderr(), status.Group[*wisski.WissKI, error]{
		PrefixString: func(item *wisski.WissKI, index int) string {
			return fmt.Sprintf("[up %q]: ", item.Slug)
		},
		PrefixAlign: true,
		Handler: func(item *wisski.WissKI, index int, writer io.Writer) (e error) {
			stack, err := item.Barrel().OpenStack()
			if err != nil {
				return fmt.Errorf("failed to open stack: %w", err)
			}
			defer errorsx.Close(stack, &e, "stack")
			return stack.Start(cmd.Context(), writer)
		},
	}, wissKIs); err != nil {
		return fmt.Errorf("failed to start instances: %w", err)
	}

	return nil
}

func (sp *systempause) stop(cmd *cobra.Command, dis *dis.Distillery) error {
	if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Shutting Down WissKIs"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}

	wissKIs, err := dis.Instances().All(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to get all instances: %w", err)
	}

	if err := status.RunErrorGroup(cmd.ErrOrStderr(), status.Group[*wisski.WissKI, error]{
		PrefixString: func(item *wisski.WissKI, index int) string {
			return fmt.Sprintf("[down %q]: ", item.Slug)
		},
		PrefixAlign: true,
		Handler: func(item *wisski.WissKI, index int, writer io.Writer) (e error) {
			stack, err := item.Barrel().OpenStack()
			if err != nil {
				return fmt.Errorf("failed to open stack: %w", err)
			}
			defer errorsx.Close(stack, &e, "stack")
			return stack.Down(cmd.Context(), writer)
		},
	}, wissKIs); err != nil {
		return fmt.Errorf("failed to shutdown instances: %w", err)
	}

	if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Shutting Down Components"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}

	if err := status.RunErrorGroup(cmd.ErrOrStderr(), status.Group[component.Installable, error]{
		PrefixString: func(item component.Installable, index int) string {
			return fmt.Sprintf("[down %q]: ", item.Name())
		},
		PrefixAlign: true,
		Handler: func(item component.Installable, index int, writer io.Writer) (e error) {
			stack, err := item.OpenStack()
			if err != nil {
				return fmt.Errorf("failed to open stack: %w", err)
			}
			defer errorsx.Close(stack, &e, "stack")
			return stack.Down(cmd.Context(), writer)
		},
	}, dis.Installable()); err != nil {
		return fmt.Errorf("failed to shutdown core instances: %w", err)
	}

	return nil
}
