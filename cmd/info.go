package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
)

// Info is then 'info' command
var Info wisski_distillery.Command = info{}

type info struct {
	Positionals struct {
		Slug string `positional-arg-name:"SLUG" required:"1-1" description:"slug of instance to show info about"`
	} `positional-args:"true"`
}

func (info) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: core.Requirements{
			NeedsDistillery: true,
		},
		Command:     "info",
		Description: "Provide information about a single repository",
	}
}

func (i info) Run(context wisski_distillery.Context) error {
	instance, err := context.Environment.Instances().WissKI(i.Positionals.Slug)
	if err != nil {
		return err
	}

	context.Printf("URL:                  %s\n", instance.URL())
	context.Printf("Base directory:       %s\n", instance.FilesystemBase)

	context.Printf("SQL Database:         %s\n", instance.SqlDatabase)
	context.Printf("SQL Username:         %s\n", instance.SqlUsername)
	context.Printf("SQL Password:         %s\n", instance.SqlPassword)

	context.Printf("GraphDB Repository:   %s\n", instance.GraphDBRepository)
	context.Printf("GraphDB Username:     %s\n", instance.GraphDBUsername)
	context.Printf("GraphDB Password:     %s\n", instance.GraphDBPassword)

	return nil
}
