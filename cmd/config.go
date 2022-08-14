package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/env"
)

// Config is the configuration command
var Config wisski_distillery.Command = config{}

type config struct {
}

func (s config) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: env.Requirements{
			NeedsConfig: true,
		},
		Command:     "config",
		Description: "Prints information about configuration",
	}
}

func (s config) Run(context wisski_distillery.Context) error {
	context.Printf("%#v", context.Environment.Config)
	return nil
}
