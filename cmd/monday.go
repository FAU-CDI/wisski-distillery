package cmd

//spellchecker:words github wisski distillery internal logging pkglib
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/pkglib/fsx"
)

// Monday is the 'monday' command.
var Monday wisski_distillery.Command = monday{}

type monday struct {
	UpdateInstances bool `short:"u" long:"update-instances" description:"fully update instances. may take a long time, and is potentially breaking"`
	SkipBackup      bool `long:"skip-backup" description:"skip making a backup. dangerous"`
	Positionals     struct {
		GraphdbZip string "positional-arg-name:\"PATH_TO_GRAPHDB_ZIP\" required:\"1-1\" description:\"path to the `graphdb.zip` file\""
	} `positional-args:"true"`
}

func (monday) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "monday",
		Description: "runs regular monday tasks",
	}
}

func (monday monday) AfterParse() error {
	isFile, err := fsx.IsRegular(monday.Positionals.GraphdbZip, false)
	if err != nil {
		return fmt.Errorf("failed to check for regular file: %w", err)
	}
	if !isFile {
		return fmt.Errorf("%q: %w", monday.Positionals.GraphdbZip, errNoGraphDBZip)
	}
	return nil
}

func (monday monday) Run(context wisski_distillery.Context) error {
	if !monday.SkipBackup {
		if err := logging.LogOperation(func() error {
			return context.Exec("backup")
		}, context.Stderr, "Running backup"); err != nil {
			return fmt.Errorf("failed to run backup: %w", err)
		}
	}

	if err := logging.LogOperation(func() error {
		return context.Exec("system_update", monday.Positionals.GraphdbZip)
	}, context.Stderr, "Running system_update"); err != nil {
		return fmt.Errorf("failed to run system_update: %w", err)
	}

	if err := logging.LogOperation(func() error {
		return context.Exec("rebuild")
	}, context.Stderr, "Running rebuild"); err != nil {
		return fmt.Errorf("failed to rebuld: %w", err)
	}

	if err := logging.LogOperation(func() error {
		return context.Exec("update_prefix_config")
	}, context.Stderr, "Running update_prefix_config"); err != nil {
		return fmt.Errorf("failed to run update_prefix_config: %w", err)
	}

	if monday.UpdateInstances {
		if err := logging.LogOperation(func() error {
			return context.Exec("blind_update")
		}, context.Stderr, "Running blind_update"); err != nil {
			return fmt.Errorf("failed to run blind_update: %w", err)
		}
	}

	if _, err := logging.LogMessage(context.Stderr, "Done, have a great week!"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	return nil
}
