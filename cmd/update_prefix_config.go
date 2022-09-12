package cmd

import (
	"io/fs"
	"os"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
)

// Cron is the 'cron' command
var UpdatePrefixConfig wisski_distillery.Command = updateprefixconfig{}

type updateprefixconfig struct{}

func (updateprefixconfig) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: core.Requirements{
			NeedsDistillery: true,
		},
		Command:     "update_prefix_config",
		Description: "Updates the prefix configuration",
	}
}

var errPrefixUpdateFailed = exit.Error{
	Message:  "Failed to update the prefix configuration: %s",
	ExitCode: exit.ExitGeneric,
}

func (upc updateprefixconfig) Run(context wisski_distillery.Context) error {
	dis := context.Environment

	instances, err := dis.AllInstances()
	if err != nil {
		return errPrefixUpdateFailed.WithMessageF(err)
	}

	resolver := dis.Resolver()
	target := resolver.ConfigPath()

	// print the configuration
	config, err := os.OpenFile(target, os.O_WRONLY, fs.ModePerm)
	if err != nil {
		return errPrefixUpdateFailed.WithMessageF(err)
	}

	// iterate over the instances and store the last value of error
	for _, instance := range instances {
		if err := logging.LogOperation(func() error {
			// read the prefix config
			data, err := instance.PrefixConfig()
			if err != nil {
				return err
			}
			context.IOStream.Printf("%s", data)

			// and write it out!
			if _, err := config.WriteString(data); err != nil {
				return err
			}

			return nil
		}, context.IOStream, "reading prefix config %s", instance.Slug); err != nil {
			return errPrefixUpdateFailed.WithMessageF(err)
		}
	}

	// and restart the resolver to apply the config!
	logging.LogMessage(context.IOStream, "restarting resolver stack")
	if err := resolver.Stack().Restart(context.IOStream); err != nil {
		return errPrefixUpdateFailed.WithMessageF(err)
	}

	return err
}
