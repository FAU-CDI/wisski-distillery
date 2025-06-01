package cmd

//spellchecker:words encoding json github wisski distillery internal goprogram exit pkglib collection
import (
	"encoding/json"
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/collection"
)

// Info is then 'info' command.
var Info wisski_distillery.Command = info{}

type info struct {
	JSON        bool `description:"print information as JSON instead of as string" long:"json" short:"j"`
	Positionals struct {
		Slug string `description:"slug of instance to show info about" positional-arg-name:"SLUG" required:"1-1"`
	} `positional-args:"true"`
}

func (info) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "info",
		Description: "provide information about a single instance",
	}
}

var errInfoFailed = exit.NewErrorWithCode("failed to get info", exit.ExitGeneric)

func (i info) Run(context wisski_distillery.Context) (err error) {
	if err := i.run(context); err != nil {
		return fmt.Errorf("%w: %w", errInfoFailed, err)
	}
	return nil
}

func (i info) run(context wisski_distillery.Context) (err error) {
	instance, err := context.Environment.Instances().WissKI(context.Context, i.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("failed to get instance: %w", err)
	}

	info, err := instance.Info().Information(context.Context, false)
	if err != nil {
		return fmt.Errorf("failed to get info: %w", err)
	}

	if i.JSON {
		if err := json.NewEncoder(context.Stdout).Encode(info); err != nil {
			return fmt.Errorf("failed to encode info as json: %w", err)
		}
		return nil
	}

	_, _ = context.Printf("Slug:                 %v\n", info.Slug)
	_, _ = context.Printf("URL:                  %v\n", info.URL)

	_, _ = context.Printf("Base directory:       %v\n", instance.FilesystemBase)

	_, _ = context.Printf("SQL Database:         %v\n", instance.SqlDatabase)
	_, _ = context.Printf("SQL Username:         %v\n", instance.SqlUsername)
	_, _ = context.Printf("SQL Password:         %v\n", instance.SqlPassword)

	_, _ = context.Printf("GraphDB Repository:   %v\n", instance.GraphDBRepository)
	_, _ = context.Printf("GraphDB Username:     %v\n", instance.GraphDBUsername)
	_, _ = context.Printf("GraphDB Password:     %v\n", instance.GraphDBPassword)

	_, _ = context.Printf("Running:              %v\n", info.Running)
	_, _ = context.Printf("Locked:               %v\n", info.Locked)
	_, _ = context.Printf("Last Rebuild:         %v\n", info.LastRebuild.String())
	_, _ = context.Printf("Last Update:          %v\n", info.LastUpdate.String())
	_, _ = context.Printf("Last Cron:            %v\n", info.LastCron.String())

	_, _ = context.Printf("Drupal Version:       %v\n", info.DrupalVersion)
	_, _ = context.Printf("Theme:                %v\n", info.Theme)

	_, _ = context.Printf("Bundles: (count %d)\n", info.Statistics.Bundles.TotalBundles)
	for _, bundle := range info.Statistics.Bundles.Bundles {
		if bundle.Count == 0 {
			continue
		}
		_, _ = context.Printf("- %s %d %v\n", bundle.Label, bundle.Count, bundle.MainBundle)
	}
	_, _ = context.Printf("Graphs: (count %d)\n", len(info.Statistics.Triplestore.Graphs))
	for _, graph := range info.Statistics.Triplestore.Graphs {
		_, _ = context.Printf("- %s %d\n", graph.URI, graph.Count)
	}

	_, _ = context.Printf("SSH Keys: (count %d)\n", len(info.SSHKeys))
	for _, key := range info.SSHKeys {
		_, _ = context.Printf("- %s\n", key)
	}

	_, _ = context.Printf("Skip Prefixes:        %v\n", info.NoPrefixes)
	_, _ = context.Printf("Prefixes: (count %d)\n", len(info.Prefixes))
	for _, prefix := range info.Prefixes {
		_, _ = context.Printf("- %s\n", prefix)
	}

	_, _ = context.Printf("Snapshots: (count %d)\n", len(info.Snapshots))
	for _, s := range info.Snapshots {
		_, _ = context.Printf("- %s (taken %s, packed %v)\n", s.Path, s.Created.String(), s.Packed)
	}

	_, _ = context.Printf("Pathbuilders: (count %d)\n", len(info.Pathbuilders))
	for name, data := range collection.IterSorted(info.Pathbuilders) {
		_, _ = context.Printf("- %s (%d bytes)\n", name, len(data))
	}

	_, _ = context.Printf("Users: (count %d)\n", len(info.Users))
	for _, user := range info.Users {
		_, _ = context.Printf("- %v\n", user)
	}

	_, _ = context.Printf("Grants: (count %d)\n", len(info.Grants))
	for _, grant := range info.Grants {
		_, _ = context.Printf("- %v\n", grant)
	}

	_, _ = context.Printf("Requirements: (count %d)\n", len(info.Requirements))
	for _, req := range info.Requirements {
		_, _ = context.Printf("- %v\n", req)
	}

	return nil
}
