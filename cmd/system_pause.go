package cmd

//spellchecker:words github wisski distillery internal component logging goprogram exit pkglib errorsx status
import (
	"fmt"
	"io"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"go.tkw01536.de/goprogram/exit"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/status"
)

// SystemPause is the 'system_pause' command.
var SystemPause wisski_distillery.Command = systempause{}

type systempause struct {
	Stop  bool `description:"stop all the components"  long:"stop"  short:"d"`
	Start bool `description:"start all the components" long:"start" short:"u"`
}

func (systempause) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "system_pause",
		Description: "stops or starts the entire distillery system",
	}
}

var errStopStartExcluded = exit.NewErrorWithCode("exactly one of `--stop` and `--start` must be provided", exit.ExitCommandArguments)

func (s systempause) AfterParse() error {
	if s.Stop == s.Start {
		return errStopStartExcluded
	}
	return nil
}

var errPauseGeneric = exit.NewErrorWithCode("unable to pause or resume system", exit.ExitGeneric)

func (sp systempause) Run(context wisski_distillery.Context) (err error) {
	if sp.Start {
		err = sp.start(context, context.Environment)
	} else {
		err = sp.stop(context, context.Environment)
	}

	if err != nil {
		return fmt.Errorf("%w: %w", errPauseGeneric, err)
	}
	return nil
}

func (sp systempause) start(context wisski_distillery.Context, dis *dis.Distillery) error {
	if _, err := logging.LogMessage(context.Stderr, "Starting Components"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}

	// find all the core stacks
	if err := status.RunErrorGroup(context.Stderr, status.Group[component.Installable, error]{
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

			return stack.Start(context.Context, writer)
		},
	}, dis.Installable()); err != nil {
		return fmt.Errorf("failed to start components: %w", err)
	}

	if _, err := logging.LogMessage(context.Stderr, "Starting Up WissKIs"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}

	// find the instances
	wissKIs, err := dis.Instances().All(context.Context)
	if err != nil {
		return fmt.Errorf("failed to get all instances: %w", err)
	}

	// shut them all down
	if err := status.RunErrorGroup(context.Stderr, status.Group[*wisski.WissKI, error]{
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

			return stack.Start(context.Context, writer)
		},
	}, wissKIs); err != nil {
		return fmt.Errorf("failed to start instances: %w", err)
	}

	return nil
}

func (sp systempause) stop(context wisski_distillery.Context, dis *dis.Distillery) error {
	if _, err := logging.LogMessage(context.Stderr, "Shutting Down WissKIs"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}

	// find the instances
	wissKIs, err := dis.Instances().All(context.Context)
	if err != nil {
		return fmt.Errorf("failed to get all instances: %w", err)
	}

	// shut them all down
	if err := status.RunErrorGroup(context.Stderr, status.Group[*wisski.WissKI, error]{
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

			return stack.Down(context.Context, writer)
		},
	}, wissKIs); err != nil {
		return fmt.Errorf("failed to shutdown instances: %w", err)
	}

	if _, err := logging.LogMessage(context.Stderr, "Shutting Down Components"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}

	// find all the core stacks
	if err := status.RunErrorGroup(context.Stderr, status.Group[component.Installable, error]{
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

			return stack.Down(context.Context, writer)
		},
	}, dis.Installable()); err != nil {
		return fmt.Errorf("failed to shutdown core instances: %w", err)
	}

	return nil
}
