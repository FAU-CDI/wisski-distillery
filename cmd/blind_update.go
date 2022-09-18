package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/tkw1536/goprogram/exit"
)

// BlindUpdate is the 'blind_update' command
var BlindUpdate wisski_distillery.Command = blindUpdate{}

type blindUpdate struct {
	Force       bool `short:"f" long:"force" description:"force running blind-update even if AutoBlindUpdate is set to false"`
	Positionals struct {
		Slug []string `positional-arg-name:"SLUG" required:"0" description:"slug of instance(s) to run blind-update in"`
	} `positional-args:"true"`
}

func (blindUpdate) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: core.Requirements{
			NeedsDistillery: true,
		},
		Command:     "blind_update",
		Description: "Runs the blind update in the provided instances",
	}
}

var errBlindUpdateFailed = exit.Error{
	Message:  "Failed to run blind update script for instance %q: exited with code %s",
	ExitCode: exit.ExitGeneric,
}

func (bu blindUpdate) Run(context wisski_distillery.Context) error {
	instances, err := context.Environment.Instances().Load(bu.Positionals.Slug...)
	if err != nil {
		return err
	}

	for _, instance := range instances {
		if !(instance.IsBlindUpdateEnabled() || bu.Force) {
			context.EPrintf("skipping instance %q\n", instance.Slug)
			continue
		}
		context.EPrintf("Updating instance %s\n", instance.Slug)

		code, err := instance.Shell(context.IOStream, "/runtime/blind_update.sh")
		if err != nil {
			return errBlindUpdateFailed.WithMessageF(instance.Slug, environment.ExecCommandError)
		}
		if code != 0 {
			return errBlindUpdateFailed.WithMessageF(instance.Slug, code)
		}
	}

	return nil
}
