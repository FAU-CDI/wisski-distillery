package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
)

// InstancePause is the 'instance_pause' command
var InstancePause wisski_distillery.Command = instancepause{}

type instancepause struct {
	Stop        bool `short:"d" long:"stop" description:"stop instance"`
	Start       bool `short:"u" long:"start" description:"start (or restart) instance"`
	Positionals struct {
		Slug string `positional-arg-name:"slug" required:"1-1" description:"name of instance to purge"`
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

var errInstancePauseWissKI = exit.Error{
	Message:  "unable to get WissKI",
	ExitCode: exit.ExitGeneric,
}

func (i instancepause) Run(context wisski_distillery.Context) error {
	instance, err := context.Environment.Instances().WissKI(context.Context, i.Positionals.Slug)
	if err != nil {
		return errInstancePauseWissKI.WrapError(err)
	}

	if i.Stop {
		return instance.Barrel().Stack().Down(context.Context, context.Stdout)
	} else {
		return instance.Barrel().Stack().Up(context.Context, context.Stdout)
	}
}
