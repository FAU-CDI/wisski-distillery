package cmd

//spellchecker:words encoding json github wisski distillery internal goprogram exit
import (
	"encoding/json"
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
)

// Info is then 'info' command.
var Status wisski_distillery.Command = cStatus{}

type cStatus struct {
	JSON bool `description:"print status as JSON instead of as string" long:"json" short:"j"`
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

var errStatusGeneric = exit.NewErrorWithCode("unable to get status", exit.ExitGeneric)

func (s cStatus) Run(context wisski_distillery.Context) error {
	status, _, err := context.Environment.Info().Status(context.Context, true)
	if err != nil {
		return fmt.Errorf("%w: %w", errStatusGeneric, err)
	}

	if s.JSON {
		err := json.NewEncoder(context.Stdout).Encode(status)
		if err != nil {
			return fmt.Errorf("%w: %w", errStatusGeneric, err)
		}
		return nil
	}

	_, _ = context.Printf("Total Instances:      %v\n", status.TotalCount)
	_, _ = context.Printf("      (running):      %v\n", status.RunningCount)
	_, _ = context.Printf("      (stopped):      %v\n", status.StoppedCount)

	_, _ = context.Printf("Backups: (count %d)\n", len(status.Backups))
	for _, s := range status.Backups {
		_, _ = context.Printf("- %s (slug %q, taken %s, packed %v)\n", s.Path, s.Slug, s.Created.String(), s.Packed)
	}

	return nil
}
