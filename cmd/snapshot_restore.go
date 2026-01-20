package cmd

//spellchecker:words bufio encoding json path filepath strings github wisski distillery internal component exporter cobra pkglib exit
import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/fsx"
)

func NewSnapshotRestoreCommand() *cobra.Command {
	impl := new(snapshotRestore)

	cmd := &cobra.Command{
		Use:     "snapshot_restore DIRECTORY SLUG",
		Short:   "restores an instance from a snapshot directory",
		Args:    cobra.ExactArgs(2),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.Yes, "yes", false, "do not ask for confirmation")

	return cmd
}

type snapshotRestore struct {
	Yes         bool
	Positionals struct {
		Slug      string // instance to restore to
		Directory string // path to the snapshot directory
	}
}

func (sr *snapshotRestore) ParseArgs(cmd *cobra.Command, args []string) error {
	sr.Positionals.Directory = args[0]
	sr.Positionals.Slug = args[1]
	return nil
}

var (
	errSnapshotRestoreNoConfirmation = exit.NewErrorWithCode("aborting after request was not confirmed. either type `yes` or pass `--yes` on the command line", cli.ExitGeneric)
	errSnapshotRestoreFailed         = exit.NewErrorWithCode("failed to restore snapshot", cli.ExitGeneric)
)

func (sr *snapshotRestore) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errSnapshotRestoreFailed, err)
	}

	if err := sr.exec(cmd, dis); err != nil {
		return fmt.Errorf("%w: %w", errSnapshotRestoreFailed, err)
	}
	return nil
}

func (sr *snapshotRestore) exec(cmd *cobra.Command, dis *dis.Distillery) error {
	if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Loading instance to restore to"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}

	// Get information about the instance to restore to.
	instance, err := dis.Instances().WissKI(cmd.Context(), sr.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("instance to restore to does not exist: %w", err)
	}

	if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Loading snapshot"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}

	snapshot, err := loadArchive(sr.Positionals.Directory)
	if err != nil {
		return fmt.Errorf("failed to load snapshot: %w", err)
	}

	// check the confirmation from the user
	if !sr.Yes {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "About to restore instance %q from %q (taken at %s). This will overwrite existing data.\n", sr.Positionals.Slug, sr.Positionals.Directory, snapshot.StartTime.Format(time.RFC3339))
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Type 'yes' to continue: ")
		reader := bufio.NewReader(cmd.InOrStdin())
		line, err := reader.ReadString('\n')
		if err != nil || strings.TrimSpace(line) != "yes" {
			return errSnapshotRestoreNoConfirmation
		}
	}

	if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Preparing restoration"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}

	// TODO: Shutdown instance to restore
	// TODO: Restore filesystem
	// TODO: Restore sql database
	// TODO: Restore triplestore
	// TODO: Start instance
	// TODO: Re-do sql config
	// TODO: Re-start instance
	// TODO: Re-do adapter configs

	// TODO: Actually implement restoration (@ai: dont do this yet)
	_ = instance // for now: to remove the unused variable error
	return exit.NewErrorWithCode("not implemented", cli.ExitGeneric)
}

func loadArchive(path string) (exporter.Snapshot, error) {
	isDirectory, err := fsx.IsDirectory(path, false)
	if err != nil {
		return exporter.Snapshot{}, fmt.Errorf("failed to check if path is a directory: %w", err)
	}
	if !isDirectory {
		return exporter.Snapshot{}, fmt.Errorf("path is not a directory: %s", path)
	}

	reportPath := filepath.Join(path, exporter.ReportMachinePath)
	reportFile, err := os.Open(reportPath)
	if err != nil {
		return exporter.Snapshot{}, fmt.Errorf("failed to open snapshot report: %w", err)
	}
	defer reportFile.Close()

	var snapshot exporter.Snapshot
	if err := json.NewDecoder(reportFile).Decode(&snapshot); err != nil {
		return exporter.Snapshot{}, fmt.Errorf("failed to decode snapshot report: %w", err)
	}

	return snapshot, nil
}
