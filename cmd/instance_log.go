package cmd

//spellchecker:words github wisski distillery internal goprogram exit pkglib errorsx
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/errorsx"
)

// InstancePause is the 'instance_log' command.
var InstanceLog wisski_distillery.Command = instanceLog{}

type instanceLog struct {
	Positionals struct {
		Slug string `description:"name of instance to follow logs for" positional-arg-name:"slug" required:"1-1"`
	} `positional-args:"true"`
}

func (instanceLog) Description() wisski_distillery.Description {
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

func (i instanceLog) Run(context wisski_distillery.Context) (e error) {
	instance, err := context.Environment.Instances().WissKI(context.Context, i.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: %w", errInstanceLogWissKI, err)
	}

	stack, err := instance.Barrel().OpenStack()
	if err != nil {
		return fmt.Errorf("%w: %w", errInstanceLogStack, err)
	}
	defer errorsx.Close(stack, &e, "stack")

	if err := stack.Attach(context.Context, context.IOStream, false); err != nil {
		return fmt.Errorf("%w: %w", errInstanceLogAttach, err)
	}
	return nil
}
