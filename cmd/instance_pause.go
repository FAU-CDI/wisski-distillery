package cmd

//spellchecker:words github wisski distillery internal goprogram exit pkglib errorsx
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/errorsx"
)

// InstancePause is the 'instance_pause' command.
var InstancePause wisski_distillery.Command = instancepause{}

type instancepause struct {
	Stop        bool `description:"stop instance"               long:"stop"  short:"d"`
	Start       bool `description:"start (or restart) instance" long:"start" short:"u"`
	Positionals struct {
		Slug string `description:"name of instance to purge" positional-arg-name:"slug" required:"1-1"`
	} `positional-args:"true"`
}

func (instancepause) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "instance_pause",
		Description: "stops or starts a single instance",
	}
}

func (i instancepause) AfterParse() error {
	if i.Stop == i.Start {
		return errStopStartExcluded
	}
	return nil
}

var (
	errInstancePauseWissKI = exit.NewErrorWithCode("unable to get WissKI", exit.ExitGeneric)
	errInstancePauseStack  = exit.NewErrorWithCode("unable to get stack", exit.ExitGeneric)
)

func (i instancepause) Run(context wisski_distillery.Context) (e error) {
	instance, err := context.Environment.Instances().WissKI(context.Context, i.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: %w", errInstancePauseWissKI, err)
	}

	stack, err := instance.Barrel().OpenStack()
	if err != nil {
		return fmt.Errorf("%w: %w", errInstancePauseStack, err)
	}
	defer errorsx.Close(stack, &e, "stack")

	if i.Stop {
		if err := stack.Down(context.Context, context.Stdout); err != nil {
			return fmt.Errorf("failed to stop instance: %w", err)
		}
	} else {
		if err := stack.Start(context.Context, context.Stdout); err != nil {
			return fmt.Errorf("failed to start instance: %w", err)
		}
	}
	return nil
}
