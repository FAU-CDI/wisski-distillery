package cmd

//spellchecker:words github wisski distillery internal logging pkglib
import (
	"errors"
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/fsx"
)

func NewMondayCommand() *cobra.Command {
	impl := new(monday)

	cmd := &cobra.Command{
		Use:     "monday GRAPHDB_ZIP",
		Short:   "runs regular monday tasks",
		Args:    cobra.ExactArgs(1),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.UpdateInstances, "update-instances", false, "fully update instances. may take a long time, and is potentially breaking")
	flags.BoolVar(&impl.SkipBackup, "skip-backup", false, "skip making a backup. dangerous")

	return cmd
}

type monday struct {
	UpdateInstances bool
	SkipBackup      bool
	Positionals     struct {
		GraphdbZip string
	}
}

func (m *monday) ParseArgs(cmd *cobra.Command, args []string) error {
	if len(args) >= 1 {
		m.Positionals.GraphdbZip = args[0]
	}

	isFile, err := fsx.IsRegular(m.Positionals.GraphdbZip, false)
	if err != nil {
		return fmt.Errorf("failed to check for regular file: %w", err)
	}
	if !isFile {
		return fmt.Errorf("%q: %w", m.Positionals.GraphdbZip, errNoGraphDBZip)
	}
	return nil
}

func (*monday) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "monday",
		Description: "runs regular monday tasks",
	}
}

var errNoGraphDBZip = errors.New("not a regular file")

func (m *monday) Exec(cmd *cobra.Command, args []string) error {
	_, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("failed to get distillery: %w", err)
	}

	if !m.SkipBackup {
		if err := logging.LogOperation(func() error {
			return execWdcli(cmd, []string{"backup"})
		}, cmd.ErrOrStderr(), "Running backup"); err != nil {
			return fmt.Errorf("failed to run backup: %w", err)
		}
	}

	if err := logging.LogOperation(func() error {
		return execWdcli(cmd, []string{"system_update", m.Positionals.GraphdbZip})
	}, cmd.ErrOrStderr(), "Running system_update"); err != nil {
		return fmt.Errorf("failed to run system_update: %w", err)
	}

	if err := logging.LogOperation(func() error {
		return execWdcli(cmd, []string{"rebuild"})
	}, cmd.ErrOrStderr(), "Running rebuild"); err != nil {
		return fmt.Errorf("failed to rebuld: %w", err)
	}

	if err := logging.LogOperation(func() error {
		return execWdcli(cmd, []string{"update_prefix_config"})
	}, cmd.ErrOrStderr(), "Running update_prefix_config"); err != nil {
		return fmt.Errorf("failed to run update_prefix_config: %w", err)
	}

	if m.UpdateInstances {
		if err := logging.LogOperation(func() error {
			return execWdcli(cmd, []string{"blind_update"})
		}, cmd.ErrOrStderr(), "Running blind_update"); err != nil {
			return fmt.Errorf("failed to run blind_update: %w", err)
		}
	}

	if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Done, have a great week!"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	return nil
}

func execWdcli(cmd *cobra.Command, args []string) error {
	wdcli := NewCommand(cmd.Context(), cli.Params{})

	wdcli.SetIn(cmd.InOrStdin())
	wdcli.SetOut(cmd.OutOrStdout())
	wdcli.SetErr(cmd.ErrOrStderr())
	wdcli.SetArgs(args)

	// despite passing the context above
	// we override flags and parameter with this!
	wdcli.SetContext(cmd.Context())

	if err := wdcli.Execute(); err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}
	return nil
}
