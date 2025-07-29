package cmd

//spellchecker:words encoding json github wisski distillery internal goprogram exit pkglib collection
import (
	"encoding/json"
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/collection"
	"go.tkw01536.de/pkglib/exit"
)

func NewInfoCommand() *cobra.Command {
	impl := new(info)

	cmd := &cobra.Command{
		Use:     "info SLUG",
		Short:   "provide information about a single instance",
		Args:    cobra.ExactArgs(1),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.JSON, "json", false, "print information as JSON instead of as string")

	return cmd
}

type info struct {
	JSON        bool
	Positionals struct {
		Slug string
	}
}

func (i *info) ParseArgs(cmd *cobra.Command, args []string) error {
	i.Positionals.Slug = args[0]
	return nil
}

func (*info) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "info",
		Description: "provide information about a single instance",
	}
}

var errInfoFailed = exit.NewErrorWithCode("failed to get info", exit.ExitGeneric)

func (i *info) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errInfoFailed, err)
	}

	if err := i.exec(cmd, dis); err != nil {
		return fmt.Errorf("%w: %w", errInfoFailed, err)
	}
	return nil
}

func (i *info) exec(cmd *cobra.Command, dis *dis.Distillery) (err error) {
	instance, err := dis.Instances().WissKI(cmd.Context(), i.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("failed to get instance: %w", err)
	}

	info, err := instance.Info().Information(cmd.Context(), false)
	if err != nil {
		return fmt.Errorf("failed to get info: %w", err)
	}

	if i.JSON {
		if err := json.NewEncoder(cmd.OutOrStdout()).Encode(info); err != nil {
			return fmt.Errorf("failed to encode info as json: %w", err)
		}
		return nil
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Slug:                 %v\n", info.Slug)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "URL:                  %v\n", info.URL)

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Base directory:       %v\n", instance.FilesystemBase)

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "SQL Database:         %v\n", instance.SqlDatabase)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "SQL Username:         %v\n", instance.SqlUsername)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "SQL Password:         %v\n", instance.SqlPassword)

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "GraphDB Repository:   %v\n", instance.GraphDBRepository)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "GraphDB Username:     %v\n", instance.GraphDBUsername)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "GraphDB Password:     %v\n", instance.GraphDBPassword)

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Running:              %v\n", info.Running)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Locked:               %v\n", info.Locked)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Last Rebuild:         %v\n", info.LastRebuild.String())
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Last Update:          %v\n", info.LastUpdate.String())
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Last Cron:            %v\n", info.LastCron.String())

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Drupal Version:       %v\n", info.DrupalVersion)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Theme:                %v\n", info.Theme)

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Bundles: (count %d)\n", info.Statistics.Bundles.TotalBundles)
	for _, bundle := range info.Statistics.Bundles.Bundles {
		if bundle.Count == 0 {
			continue
		}
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "- %s %d %v\n", bundle.Label, bundle.Count, bundle.MainBundle)
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Graphs: (count %d)\n", len(info.Statistics.Triplestore.Graphs))
	for _, graph := range info.Statistics.Triplestore.Graphs {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "- %s %d\n", graph.URI, graph.Count)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "SSH Keys: (count %d)\n", len(info.SSHKeys))
	for _, key := range info.SSHKeys {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", key)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Skip Prefixes:        %v\n", info.NoPrefixes)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Prefixes: (count %d)\n", len(info.Prefixes))
	for _, prefix := range info.Prefixes {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", prefix)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Snapshots: (count %d)\n", len(info.Snapshots))
	for _, s := range info.Snapshots {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "- %s (taken %s, packed %v)\n", s.Path, s.Created.String(), s.Packed)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Pathbuilders: (count %d)\n", len(info.Pathbuilders))
	for name, data := range collection.IterSorted(info.Pathbuilders) {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "- %s (%d bytes)\n", name, len(data))
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Users: (count %d)\n", len(info.Users))
	for _, user := range info.Users {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "- %v\n", user)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Grants: (count %d)\n", len(info.Grants))
	for _, grant := range info.Grants {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "- %v\n", grant)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Requirements: (count %d)\n", len(info.Requirements))
	for _, req := range info.Requirements {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "- %v\n", req)
	}

	return nil
}
