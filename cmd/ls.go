package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/env"
)

// Ls is the 'ls' command
var Ls wisski_distillery.Command = ls{}

type ls struct {
	Positionals struct {
		Slug []string `positional-arg-name:"SLUG" required:"0" description:"slug(s) of instance(s) to list"`
	} `positional-args:"true"`
}

func (ls) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: env.Requirements{
			NeedsDistillery: true,
		},
		Command:     "ls",
		Description: "Lists WissKI instances",
	}
}

func (l ls) Run(context wisski_distillery.Context) error {
	instances, err := context.Environment.Instances(l.Positionals.Slug...)
	if err != nil {
		return err
	}

	for _, instance := range instances {
		context.Println(instance.Slug)
	}

	return nil
}
