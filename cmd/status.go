package cmd

import (
	"encoding/json"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
)

// Info is then 'info' command
var Status wisski_distillery.Command = cStatus{}

type cStatus struct {
	JSON bool `short:"j" long:"json" description:"print status as JSON instead of as string"`
}

func (cStatus) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "status",
		Description: "provide information about the distillery as a whole",
	}
}

func (s cStatus) Run(context wisski_distillery.Context) error {
	status, _, err := context.Environment.Info().Status(context.Context, true)
	if err != nil {
		return err
	}

	if s.JSON {
		json.NewEncoder(context.Stdout).Encode(status)
		return nil
	}

	context.Printf("Total Instances:      %v\n", status.TotalCount)
	context.Printf("      (running):      %v\n", status.RunningCount)
	context.Printf("      (stopped):      %v\n", status.StoppedCount)

	context.Printf("Backups: (count %d)\n", len(status.Backups))
	for _, s := range status.Backups {
		context.Printf("- %s (slug %q, taken %s, packed %v)\n", s.Path, s.Slug, s.Created.String(), s.Packed)
	}

	return nil
}
