package cmd

//spellchecker:words encoding json github wisski distillery internal goprogram exit pkglib collection
import (
	"encoding/json"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/pkg/errwrap"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/collection"
)

// Info is then 'info' command
var Info wisski_distillery.Command = info{}

type info struct {
	JSON        bool `short:"j" long:"json" description:"print information as JSON instead of as string"`
	Positionals struct {
		Slug string `positional-arg-name:"SLUG" required:"1-1" description:"slug of instance to show info about"`
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

var errInfoFailed = exit.Error{
	Message:  "failed to get info",
	ExitCode: exit.ExitGeneric,
}

func (i info) Run(context wisski_distillery.Context) (err error) {
	defer errwrap.DeferWrap(errInfoFailed, &err)

	instance, err := context.Environment.Instances().WissKI(context.Context, i.Positionals.Slug)
	if err != nil {
		return err
	}

	info, err := instance.Info().Information(context.Context, false)
	if err != nil {
		return err
	}

	if i.JSON {
		return json.NewEncoder(context.Stdout).Encode(info)
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
	context.Printf("Locked:               %v\n", info.Locked)
	context.Printf("Last Rebuild:         %v\n", info.LastRebuild.String())
	context.Printf("Last Update:          %v\n", info.LastUpdate.String())
	context.Printf("Last Cron:            %v\n", info.LastCron.String())

	context.Printf("Drupal Version:       %v\n", info.DrupalVersion)
	context.Printf("Theme:                %v\n", info.Theme)

	context.Printf("Bundles: (count %d)\n", info.Statistics.Bundles.TotalBundles)
	for _, bundle := range info.Statistics.Bundles.Bundles {
		if bundle.Count == 0 {
			continue
		}
		context.Printf("- %s %d %v\n", bundle.Label, bundle.Count, bundle.MainBundle)
	}
	context.Printf("Graphs: (count %d)\n", len(info.Statistics.Triplestore.Graphs))
	for _, graph := range info.Statistics.Triplestore.Graphs {
		context.Printf("- %s %d\n", graph.URI, graph.Count)
	}

	context.Printf("SSH Keys: (count %d)\n", len(info.SSHKeys))
	for _, key := range info.SSHKeys {
		context.Printf("- %s\n", key)
	}

	context.Printf("Skip Prefixes:        %v\n", info.NoPrefixes)
	context.Printf("Prefixes: (count %d)\n", len(info.Prefixes))
	for _, prefix := range info.Prefixes {
		context.Printf("- %s\n", prefix)
	}

	context.Printf("Snapshots: (count %d)\n", len(info.Snapshots))
	for _, s := range info.Snapshots {
		context.Printf("- %s (taken %s, packed %v)\n", s.Path, s.Created.String(), s.Packed)
	}

	context.Printf("Pathbuilders: (count %d)\n", len(info.Pathbuilders))
	for name, data := range collection.IterSorted(info.Pathbuilders) {
		context.Printf("- %s (%d bytes)\n", name, len(data))
	}

	context.Printf("Users: (count %d)\n", len(info.Users))
	for _, user := range info.Users {
		context.Printf("- %v\n", user)
	}

	context.Printf("Grants: (count %d)\n", len(info.Grants))
	for _, grant := range info.Grants {
		context.Printf("- %v\n", grant)
	}

	context.Printf("Requirements: (count %d)\n", len(info.Requirements))
	for _, req := range info.Requirements {
		context.Printf("- %v\n", req)
	}

	return nil
}
