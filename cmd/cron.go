package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/env"
	"github.com/FAU-CDI/wisski-distillery/internal/logging"
	"github.com/tkw1536/goprogram/exit"
)

// Cron is the 'cron' command
var Cron wisski_distillery.Command = cron{}

type cron struct {
	Positionals struct {
		Slug []string `positional-arg-name:"SLUG" required:"0" description:"slug of instance(s) to run cron in"`
	} `positional-args:"true"`
}

func (cron) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: env.Requirements{
			NeedsDistillery: true,
		},
		Command:     "cron",
		Description: "Runs the cron script for several instances",
	}
}

var errCronFailed = exit.Error{
	Message:  "Failed to run cron script for instance %q: exited with code %s",
	ExitCode: exit.ExitGeneric,
}

func (cr cron) Run(context wisski_distillery.Context) error {
	instances, err := context.Environment.Instances(cr.Positionals.Slug...)
	if err != nil {
		return err
	}

	// iterate over the instances and store the last value of error
	for _, instance := range instances {
		logging.LogOperation(func() error {
			code, err := instance.Shell(context.IOStream, "/utils/cron.sh")
			if err != nil {
				context.EPrintln(err)
			}
			if code != 0 {
				// keep going, because we want to run as many crons as possible
				err = errBlindUpdateFailed.WithMessageF(instance.Slug, code)
				context.EPrintln(err)
			}

			return nil
		}, context.IOStream, "running cron for instance %s", instance.Slug)
	}

	return err
}
