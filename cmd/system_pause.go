package cmd

//spellchecker:words github wisski distillery internal component logging goprogram exit pkglib status
import (
	"fmt"
	"io"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/status"
)

// SystemPause is the 'system_pause' command
var SystemPause wisski_distillery.Command = systempause{}

type systempause struct {
	Stop  bool `short:"d" long:"stop" description:"stop all the components"`
	Start bool `short:"u" long:"start" description:"start all the components"`
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

var errStopStartExcluded = exit.Error{
	Message:  "exactly one of `--stop` and `--start` must be provided",
	ExitCode: exit.ExitCommandArguments,
}

func (s systempause) AfterParse() error {
	if s.Stop == s.Start {
		return errStopStartExcluded
	}
	return nil
}

var errPauseGeneric = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to pause or resume system",
}

func (sp systempause) Run(context wisski_distillery.Context) (err error) {
	defer errPauseGeneric.DeferWrap(&err)

	if sp.Start {
		return sp.start(context, context.Environment)
	} else {
		return sp.stop(context, context.Environment)
	}
}

func (sp systempause) start(context wisski_distillery.Context, dis *dis.Distillery) error {
	logging.LogMessage(context.Stderr, "Starting Components")

	// find all the core stacks
	if err := status.RunErrorGroup(context.Stderr, status.Group[component.Installable, error]{
		PrefixString: func(item component.Installable, index int) string {
			return fmt.Sprintf("[up %q]: ", item.Name())
		},
		PrefixAlign: true,

		Handler: func(item component.Installable, index int, writer io.Writer) error {
			return item.Stack().Up(context.Context, writer)
		},
	}, dis.Installable()); err != nil {
		return err
	}

	logging.LogMessage(context.Stderr, "Starting Up WissKIs")

	// find the instances
	wissKIs, err := dis.Instances().All(context.Context)
	if err != nil {
		return err
	}

	// shut them all down
	if err := status.RunErrorGroup(context.Stderr, status.Group[*wisski.WissKI, error]{
		PrefixString: func(item *wisski.WissKI, index int) string {
			return fmt.Sprintf("[up %q]: ", item.Slug)
		},
		PrefixAlign: true,

		Handler: func(item *wisski.WissKI, index int, writer io.Writer) error {
			return item.Barrel().Stack().Up(context.Context, writer)
		},
	}, wissKIs); err != nil {
		return err
	}

	return nil
}

func (sp systempause) stop(context wisski_distillery.Context, dis *dis.Distillery) error {
	logging.LogMessage(context.Stderr, "Shutting Down WissKIs")

	// find the instances
	wissKIs, err := dis.Instances().All(context.Context)
	if err != nil {
		return err
	}

	// shut them all down
	if err := status.RunErrorGroup(context.Stderr, status.Group[*wisski.WissKI, error]{
		PrefixString: func(item *wisski.WissKI, index int) string {
			return fmt.Sprintf("[down %q]: ", item.Slug)
		},
		PrefixAlign: true,

		Handler: func(item *wisski.WissKI, index int, writer io.Writer) error {
			return item.Barrel().Stack().Down(context.Context, writer)
		},
	}, wissKIs); err != nil {
		return err
	}

	logging.LogMessage(context.Stderr, "Shutting Down Components")

	// find all the core stacks
	if err := status.RunErrorGroup(context.Stderr, status.Group[component.Installable, error]{
		PrefixString: func(item component.Installable, index int) string {
			return fmt.Sprintf("[down %q]: ", item.Name())
		},
		PrefixAlign: true,

		Handler: func(item component.Installable, index int, writer io.Writer) error {
			return item.Stack().Down(context.Context, writer)
		},
	}, dis.Installable()); err != nil {
		return err
	}

	return nil
}
