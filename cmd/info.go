package cmd

import (
	"encoding/json"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/pkg/slicesx"
)

// Info is then 'info' command
var Info wisski_distillery.Command = info{}

type info struct {
	JSON        bool `short:"j" long:"json" description:"Print information as JSON instead of as string"`
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

	info, err := instance.Info(false)
	if err != nil {
		return err
	}

	if i.JSON {
		json.NewEncoder(context.Stdout).Encode(info)
		return nil
	}

	context.Printf("Slug:                 %v\n", info.Slug)
	context.Printf("URL:                  %v\n", info.URL)

	context.Printf("Base directory:       %v\n", instance.FilesystemBase)

	context.Printf("SQL Database:         %v\n", instance.SqlDatabase)
	context.Printf("SQL Username:         %v\n", instance.SqlUsername)
	context.Printf("SQL Password:         %v\n", instance.SqlPassword)

	context.Printf("GraphDB Repository:   %v\n", instance.GraphDBRepository)
	context.Printf("GraphDB Username:     %v\n", instance.GraphDBUsername)
	context.Printf("GraphDB Password:     %v\n", instance.GraphDBPassword)

	context.Printf("Running:              %v\n", info.Running)
	context.Printf("Last Rebuild:         %v\n", info.LastRebuild.String())

	context.Printf("Prefixes: (count %d)\n", len(info.Prefixes))
	for _, prefix := range info.Prefixes {
		context.Printf("- %s\n", prefix)
	}

	context.Printf("Snapshots: (count %d)\n", len(info.Snapshots))
	for _, s := range info.Snapshots {
		context.Printf("- %s (taken %s, packed %v)\n", s.Path, s.Created.String(), s.Packed)
	}

	context.Printf("Pathbuilders: (count %d)\n", len(info.Pathbuilders))
	slicesx.ForSorted(info.Pathbuilders, func(name, data string) {
		context.Printf("- %s (%d bytes)\n", name, len(data))
	})

	return nil
}
