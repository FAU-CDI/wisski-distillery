package cmd

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
	"github.com/tkw1536/goprogram/status"
	"github.com/tkw1536/goprogram/stream"
)

// SystemPause is the 'system_pause' command
var SystemPause wisski_distillery.Command = systempause{}

type systempause struct {
	Stop  bool `short:"d" long:"stop" description:"Stop all the components"`
	Start bool `short:"u" long:"start" description:"Start all the components"`
}

func (systempause) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "system_pause",
		Description: "Stops or Starts the entire WissKI Distillery system",
	}
}

var errStopStartExcluded = exit.Error{
	Message:  "Exactly one of `--stop` and `--start` must be provied",
	ExitCode: exit.ExitCommandArguments,
}

func (s systempause) AfterParse() error {
	if s.Stop == s.Start {
		return errStopStartExcluded
	}
	return nil
}

func (sp systempause) Run(context wisski_distillery.Context) error {
	if sp.Start {
		return sp.start(context, context.Environment)
	} else {
		return sp.stop(context, context.Environment)
	}
}

func (sp systempause) start(context wisski_distillery.Context, dis *dis.Distillery) error {
	logging.LogMessage(context.IOStream, "Starting Components")

	// find all the core stacks
	if err := status.RunErrorGroup(context.Stdout, status.Group[component.Installable, error]{
		PrefixString: func(item component.Installable, index int) string {
			return fmt.Sprintf("[up %q]: ", item.Name())
		},
		PrefixAlign: true,

		Handler: func(item component.Installable, index int, writer io.Writer) error {
			io := stream.NewIOStream(writer, writer, stream.Null, 0)
			return item.Stack(context.Environment.Environment).Up(io)
		},
	}, dis.Installable()); err != nil {
		return err
	}

	logging.LogMessage(context.IOStream, "Starting Up WissKIs")

	// find the instances
	wissKIs, err := dis.Instances().All()
	if err != nil {
		return err
	}

	// shut them all down
	if err := status.RunErrorGroup(context.Stdout, status.Group[*wisski.WissKI, error]{
		PrefixString: func(item *wisski.WissKI, index int) string {
			return fmt.Sprintf("[up %q]: ", item.Slug)
		},
		PrefixAlign: true,

		Handler: func(item *wisski.WissKI, index int, writer io.Writer) error {
			io := stream.NewIOStream(writer, writer, stream.Null, 0)
			return item.Barrel().Stack().Up(io)
		},
	}, wissKIs); err != nil {
		return err
	}

	return nil
}

func (sp systempause) stop(context wisski_distillery.Context, dis *dis.Distillery) error {
	logging.LogMessage(context.IOStream, "Shutting Down WissKIs")

	// find the instances
	wissKIs, err := dis.Instances().All()
	if err != nil {
		return err
	}

	// shut them all down
	if err := status.RunErrorGroup(context.Stdout, status.Group[*wisski.WissKI, error]{
		PrefixString: func(item *wisski.WissKI, index int) string {
			return fmt.Sprintf("[down %q]: ", item.Slug)
		},
		PrefixAlign: true,

		Handler: func(item *wisski.WissKI, index int, writer io.Writer) error {
			io := stream.NewIOStream(writer, writer, stream.Null, 0)
			return item.Barrel().Stack().Down(io)
		},
	}, wissKIs); err != nil {
		return err
	}

	logging.LogMessage(context.IOStream, "Shutting Down Components")

	// find all the core stacks
	if err := status.RunErrorGroup(context.Stdout, status.Group[component.Installable, error]{
		PrefixString: func(item component.Installable, index int) string {
			return fmt.Sprintf("[down %q]: ", item.Name())
		},
		PrefixAlign: true,

		Handler: func(item component.Installable, index int, writer io.Writer) error {
			io := stream.NewIOStream(writer, writer, stream.Null, 0)
			return item.Stack(context.Environment.Environment).Down(io)
		},
	}, dis.Installable()); err != nil {
		return err
	}

	return nil
}

/*
	handleStack := func(io stream.IOStream, stack component.StackWithResources) error {
		if sp.Start {
			return stack.Up(io)
		} else {
			return stack.Down(io)
		}
	}

	logging.LogMessage(context.IOStream, "Iterating over components")
	if err := status.RunErrorGroup(context.Stdout, status.Group[component.Installable, error]{
		PrefixString: func(item component.Installable, index int) string {
			return fmt.Sprintf("[%s %q]: ", verb, item.Name())
		},
		PrefixAlign: true,

		Handler: func(item component.Installable, index int, writer io.Writer) error {
			io := stream.NewIOStream(writer, writer, stream.Null, 0)
			stack := item.Stack(context.Environment.Environment)

			return handleStack(io, stack)
		},
	}, dis.Installable()); err != nil {
		return err
	}

	logging.LogMessage(context.IOStream, "Shutting Down WissKIs")

	// find the instances
	wissKIs, err := dis.Instances().All()
	if err != nil {
		return err
	}

	// and do the actual rebuild
	if err := status.StreamGroup(context.IOStream, rb.Parallel, func(instance *wisski.WissKI, io stream.IOStream) error {
		return instance.Barrel().Build(io, true)
	}, wissKIs, status.SmartMessage(func(item *wisski.WissKI) string {
		return fmt.Sprintf("rebuild %q", item.Slug)
	})); err != nil {
		return err
	}
}
*/
