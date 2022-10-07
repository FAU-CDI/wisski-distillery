package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
)

// Config is the configuration command
var Config wisski_distillery.Command = cfg{}

type cfg struct{}

func (c cfg) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: core.Requirements{
			NeedsDistillery: true,
		},
		Command:     "config",
		Description: "Prints information about configuration",
	}
}

func (cfg) Run(context wisski_distillery.Context) error {
	context.Printf("%#v", context.Environment.Config)
	return nil
}
