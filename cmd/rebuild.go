package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
)

// Cron is the 'cron' command
var Rebuild wisski_distillery.Command = rebuild{}

type rebuild struct {
	Positionals struct {
		Slug []string `positional-arg-name:"SLUG" required:"0" description:"slug of instance(s) to run rebuild"`
	} `positional-args:"true"`
}

func (rebuild) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: core.Requirements{
			NeedsDistillery: true,
		},
		Command:     "rebuild",
		Description: "Runs the rebuild script for several instances",
	}
}

var errRebuildFailed = exit.Error{
	Message:  "Failed to run rebuild script for instance %q: exited with code %s",
	ExitCode: exit.ExitGeneric,
}

func (rb rebuild) Run(context wisski_distillery.Context) error {
	instances, err := context.Environment.Instances().Load(rb.Positionals.Slug...)
	if err != nil {
		return err
	}

	// iterate over the instances and store the last value of error
	var globalErr error
	for _, instance := range instances {
		logging.LogOperation(func() error {
			s := instance.Barrel()
			if err := logging.LogOperation(func() error {
				return s.Install(context.IOStream, component.InstallationContext{})
			}, context.IOStream, "Installing docker stack"); err != nil {
				globalErr = err
				return err
			}

			if err := logging.LogOperation(func() error {
				return s.Update(context.IOStream, true)
			}, context.IOStream, "Updating docker stack"); err != nil {
				globalErr = err
				return err
			}

			return nil
		}, context.IOStream, "Rebuilding instance %s", instance.Slug)
	}

	return globalErr
}
