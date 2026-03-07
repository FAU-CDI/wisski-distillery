package cmd

//spellchecker:words bufio encoding json path filepath strings github wisski distillery internal component exporter cobra pkglib exit
import (
	"bufio"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/fsx"
	"go.tkw01536.de/pkglib/status"
	"go.tkw01536.de/pkglib/stream"
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

	checkResult, err := readArchiveParts(sr.Positionals.Directory, snapshot)
	if err != nil {
		return fmt.Errorf("snapshot is not suitable for restoration: %w", err)
	}
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "================================================\n"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Instance: %s\n", instance.FilesystemBase); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Snapshot: %s\n", sr.Positionals.Directory); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "================================================\n"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Data:        %s\n", checkResult.DataPath); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "SQL:         %s\n", checkResult.SQLFilePath); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Triplestore: %s\n", checkResult.TSFilePath); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "================================================\n"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
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

	{
		if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Shutting down instance"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		if err := shutdownInstance(cmd, instance); err != nil {
			return fmt.Errorf("failed to shutdown instance: %w", err)
		}
	}

	// Data
	{
		if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Restoring data directory"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		oldDataDirectory := filepath.Join(instance.FilesystemBase, "data")

		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Old data directory: %s\n", oldDataDirectory); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "New data directory: %s\n", checkResult.DataPath); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		if err := checkResult.restoreDataDirectory(cmd, oldDataDirectory); err != nil {
			return fmt.Errorf("failed to restore data directory: %w", err)
		}

		// TODO: chown the data directory to the www-data user

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), "Data directory restored."); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
	}

	// Triplestore
	{
		if err := logging.LogOperation(func() error {
			return checkResult.restoreTriplestore(cmd, instance)
		}, cmd.ErrOrStderr(), "Restoring triplestore"); err != nil {
			return fmt.Errorf("failed to restore triplestore: %w", err)
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), "Triplestore restored."); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
	}

	// SQL
	{
		if err := logging.LogOperation(func() error {
			return checkResult.restoreSQL(cmd, instance)
		}, cmd.ErrOrStderr(), "Restoring SQL"); err != nil {
			return fmt.Errorf("failed to restore SQL database: %w", err)
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), "SQL restored."); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
	}

	// Restart instance
	// TODO: Restart in dummy mode!
	{
		if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Re-Starting instance"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		if err := startInstance(cmd, instance); err != nil {
			return fmt.Errorf("failed to restart instance: %w", err)
		}
	}

	// Re-set permissions
	{
		if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Re-setting permissions and ownership"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		barrel := instance.Barrel()
		for _, script := range [][]string{
			{"chown", "-R", "www-data:www-data", "/var/www/data/project/"},
		} {
			if err := barrel.BashScriptAs(cmd.Context(), "root", stream.NonInteractive(cmd.OutOrStdout()), script...); err != nil {
				return fmt.Errorf("failed to reset permissions: %w", err)
			}
		}
	}

	// Re-create SQL config
	{
		if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Re-Creating SQL config"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		if err := checkResult.restoreSQLConfig(cmd, instance); err != nil {
			return fmt.Errorf("failed to restart instance: %w", err)
		}
	}

	// Restart instance
	{
		if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Re-Starting instance"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		if err := startInstance(cmd, instance); err != nil {
			return fmt.Errorf("failed to restart instance: %w", err)
		}
	}

	// Re-create Adapter
	{
		if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Re-Creating Adapter"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		if _, err := instance.Adapters().SetAdapter(cmd.Context(), nil, instance.Adapters().DefaultAdapter()); err != nil {
			return fmt.Errorf("failed to restore adapter: %w", err)
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), "Adapter re-created."); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
	}

	// Re-build settings
	{
		if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Re-Applying settings"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		if err := instance.SystemManager().Apply(cmd.Context(), cmd.OutOrStdout(), instance.System); err != nil {
			return fmt.Errorf("failed to apply settings: %w", err)
		}
	}

	// Wait for SQL
	{
		if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Waiting for SQL"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		if err := waitSQL(cmd, instance); err != nil {
			return fmt.Errorf("failed to wait for SQL: %w", err)
		}
	}

	// clear cache
	{
		if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Clearing cache"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		if err := instance.Drush().Exec(cmd.Context(), cmd.OutOrStdout(), "cr"); err != nil {
			return fmt.Errorf("failed to run drush cr: %w", err)
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), "Cache cleared."); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
	}

	// and do a final restart for good measure ...
	{
		if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Re-Starting instance"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		if err := startInstance(cmd, instance); err != nil {
			return fmt.Errorf("failed to restart instance: %w", err)
		}
	}

	if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Instance should be restored."); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	return nil
}

var (
	errFailedToCheckDirectory = errors.New("failed to check if path is a directory")
	errPathNotDirectory       = errors.New("path is not a directory")
	errFailedToOpenReport     = errors.New("failed to open snapshot report")
	errFailedToDecodeReport   = errors.New("failed to decode snapshot report")

	errSnapshotMissingDataPart        = errors.New("data part is not present in the snapshot")
	errSnapshotMissingSQLPart         = errors.New("sql part missing from snapshot manifest")
	errSnapshotMissingTriplestorePart = errors.New("triplestore part missing from snapshot manifest")

	errDataFolderNotDirectory        = errors.New("data folder is not a directory")
	errSQLDataNotFoundInSnapshot     = errors.New("sql data not found in snapshot")
	errTriplestoreDataNotRegularFile = errors.New("triplestore data not a regular file")

	errFailedToOpenStack                    = errors.New("failed to open stack")
	errFailedToCreateTemporaryDirectory     = errors.New("failed to create temporary directory")
	errFailedToComputeRelativePath          = errors.New("failed to compute relative path")
	errFailedToGetDirectoryInfo             = errors.New("failed to get directory info")
	errFailedToCreateDirectory              = errors.New("failed to create directory")
	errFailedToOpenFile                     = errors.New("failed to open file")
	errFailedToCopyFile                     = errors.New("failed to copy file")
	errFailedToCopyDirectory                = errors.New("failed to copy directory")
	errFailedToRemoveOldDirectory           = errors.New("failed to remove old directory")
	errFailedToMoveRestoredDirectoryInPlace = errors.New("failed to move restored directory into place")

	errFailedToLogMessage                 = errors.New("failed to log message")
	errFailedToPurgeTriplestoreData       = errors.New("failed to purge triplestore data")
	errFailedToProvisionTriplestore       = errors.New("failed to provision triplestore")
	errFailedToOpenTriplestoreBackup      = errors.New("failed to open triplestore backup")
	errFailedToRestoreTriplestoreContents = errors.New("failed to restore triplestore contents")

	errFailedToPurgeSQLDatabase     = errors.New("failed to purge SQL database")
	errFailedToProvisionSQLDatabase = errors.New("failed to provision SQL database")
	errFailedToOpenSQLBackup        = errors.New("failed to open SQL backup")
	errFailedToRestoreSQLContents   = errors.New("failed to restore SQL contents")
	errFailedToRestoreSQLConfig     = errors.New("failed to restore SQL config")

	errFailedToLstat         = errors.New("failed to lstat")
	errFailedToWalkDir       = errors.New("failed to walk dir")
	errFailedToCreateSymlink = errors.New("failed to create symlink")

	errFailedToShutdownInstance = errors.New("failed to shutdown instance")
	errFailedToStartInstance    = errors.New("failed to start instance")
	errFailedToWaitSQL          = errors.New("failed to wait for SQL")
)

func loadArchive(path string) (s exporter.Snapshot, e error) {
	isDirectory, err := fsx.IsDirectory(path, false)
	if err != nil {
		return exporter.Snapshot{}, fmt.Errorf("%w: %w", errFailedToCheckDirectory, err)
	}
	if !isDirectory {
		return exporter.Snapshot{}, &fs.PathError{Op: "loadArchive", Path: path, Err: errPathNotDirectory}
	}

	reportPath := filepath.Join(path, exporter.ReportMachinePath)
	reportFile, err := os.Open(reportPath) // #nosec G304 -- intended
	if err != nil {
		return exporter.Snapshot{}, fmt.Errorf("%w: %w", errFailedToOpenReport, &fs.PathError{Op: "open", Path: reportPath, Err: err})
	}
	defer errorsx.Close(reportFile, &e, "report file")

	var snapshot exporter.Snapshot
	if err := json.NewDecoder(reportFile).Decode(&snapshot); err != nil {
		return exporter.Snapshot{}, fmt.Errorf("%w: %w", errFailedToDecodeReport, err)
	}

	return snapshot, nil
}

type archiveParts struct {
	DataPath string

	SQLFilePath string

	TSFilePath string
}

func readArchiveParts(path string, archive exporter.Snapshot) (parts archiveParts, err error) {
	{
		if !slices.Contains(archive.Description.Parts, "data") {
			return archiveParts{}, errSnapshotMissingDataPart
		}

		parts.DataPath = filepath.Join(path, "data", "data")
		if isDirectory, err := fsx.IsDirectory(parts.DataPath, false); !isDirectory {
			return archiveParts{}, fmt.Errorf("%w: %w", errDataFolderNotDirectory, cmp.Or(err, fs.ErrNotExist))
		}
	}

	{
		local, err := findSQLPath(archive)
		if err != nil {
			return archiveParts{}, err
		}

		parts.SQLFilePath = filepath.Join(path, local)
		if isFile, err := fsx.IsRegular(parts.SQLFilePath, false); !isFile {
			return archiveParts{}, fmt.Errorf("%w: %s: %w", errSQLDataNotFoundInSnapshot, parts.SQLFilePath, cmp.Or(err, fs.ErrNotExist))
		}
	}

	{
		if !slices.Contains(archive.Description.Parts, "triplestore") {
			return archiveParts{}, errSnapshotMissingTriplestorePart
		}

		parts.TSFilePath = filepath.Join(path, "triplestore", archive.Instance.GraphDBRepository+".nq")
		if isFile, err := fsx.IsRegular(parts.TSFilePath, false); !isFile {
			return archiveParts{}, fmt.Errorf("%w: %w", errTriplestoreDataNotRegularFile, cmp.Or(err, fs.ErrNotExist))
		}
	}

	return parts, nil
}

func findSQLPath(archive exporter.Snapshot) (string, error) {
	hasSQLPart := false
	for _, part := range archive.Description.Parts {
		if part != "sql" {
			continue
		}
		hasSQLPart = true
	}
	if !hasSQLPart {
		return "", errSnapshotMissingSQLPart
	}

	var candidates []string
	for _, path := range archive.Manifest {
		// TODO: This is a very ugly search of the manifest
		// But it's good enough for now.
		if !strings.HasPrefix(path, "sql/") || !strings.HasSuffix(path, ".sql") {
			continue
		}
		candidates = append(candidates, path)
	}

	if len(candidates) == 0 {
		return "", fmt.Errorf("%w: No sql files found in manifest", errSQLDataNotFoundInSnapshot)
	}
	if len(candidates) > 1 {
		return "", fmt.Errorf("%w: Multiple sql files found in manifest", errSQLDataNotFoundInSnapshot)
	}
	return candidates[0], nil
}

func shutdownInstance(cmd *cobra.Command, instance *wisski.WissKI) (e error) {
	stack, err := instance.Barrel().OpenStack()
	if err != nil {
		return fmt.Errorf("%w: %w", errFailedToOpenStack, err)
	}
	defer errorsx.Close(stack, &e, "stack")

	if err := stack.Down(cmd.Context(), cmd.OutOrStdout()); err != nil {
		return fmt.Errorf("%w: %w", errFailedToShutdownInstance, err)
	}
	return nil
}

func startInstance(cmd *cobra.Command, instance *wisski.WissKI) (e error) {
	stack, err := instance.Barrel().OpenStack()
	if err != nil {
		return fmt.Errorf("%w: %w", errFailedToOpenStack, err)
	}
	defer errorsx.Close(stack, &e, "stack")
	if err := stack.Start(cmd.Context(), cmd.OutOrStdout()); err != nil {
		return fmt.Errorf("%w: %w", errFailedToStartInstance, err)
	}

	return waitSQL(cmd, instance)
}

func waitSQL(cmd *cobra.Command, instance *wisski.WissKI) (e error) {
	// wait for the sql to be up
	if err := instance.BoundSQL().Impl.StartAndWait(cmd.Context(), cmd.OutOrStdout()); err != nil {
		return fmt.Errorf("%w: %w", errFailedToWaitSQL, err)
	}
	return nil
}

func (parts archiveParts) restoreDataDirectory(cmd *cobra.Command, old string) (e error) {
	fresh := parts.DataPath

	st := status.NewWithCompat(cmd.ErrOrStderr(), 1)
	st.Start()
	defer st.Set(0, "")
	defer st.Stop()

	// Create a temporary directory next to the old directory
	oldParent := filepath.Dir(old)
	tempDir, err := os.MkdirTemp(oldParent, ".restore-*")
	if err != nil {
		return fmt.Errorf("%w: %w", errFailedToCreateTemporaryDirectory, err)
	}

	// Clean up temporary directory on failure
	defer func() {
		if e != nil {
			if _, err := logging.LogMessage(cmd.ErrOrStderr(), "failed to restore, cleaning up"); err != nil {
				e = errorsx.Combine(e, fmt.Errorf("%w: %w", errFailedToLogMessage, err))
			}
			err := os.RemoveAll(tempDir)
			e = errorsx.Combine(e, err)
		}
	}()

	// Copy the new directory contents to the temporary directory
	if err := filepath.WalkDir(fresh, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Compute relative path and destination
		relPath, err := filepath.Rel(fresh, path)
		if err != nil {
			return fmt.Errorf("%w: %w", errFailedToComputeRelativePath, err)
		}
		destPath := filepath.Join(tempDir, relPath)

		st.Set(0, relPath)

		if d.IsDir() {
			// Create directory with same permissions
			info, err := d.Info()
			if err != nil {
				return fmt.Errorf("%w: %w", errFailedToGetDirectoryInfo, err)
			}
			if err := os.MkdirAll(destPath, info.Mode().Perm()); err != nil {
				return fmt.Errorf("%w: %w", errFailedToCreateDirectory, &fs.PathError{Op: "mkdirall", Path: destPath, Err: err})
			}
			return nil
		}

		// Copy regular file
		if err := copyFile(path, destPath); err != nil {
			return fmt.Errorf("%w: %w", errFailedToCopyFile, &fs.PathError{Op: "copy", Path: relPath, Err: err})
		}

		return nil
	}); err != nil {
		return fmt.Errorf("%w: %w", errFailedToCopyDirectory, err)
	}

	// Remove the old directory
	if err := os.RemoveAll(old); err != nil {
		return fmt.Errorf("%w: %w", errFailedToRemoveOldDirectory, err)
	}

	// Move the temporary directory to the old directory's place
	if err := os.Rename(tempDir, old); err != nil {
		return fmt.Errorf("%w: %w", errFailedToMoveRestoredDirectoryInPlace, err)
	}

	return nil
}

// Copies a file from src to dst.
// If it is a symlink, it is copied as a symlink.
func copyFile(src, dst string) (e error) {
	srcInfo, err := os.Lstat(src)
	if err != nil {
		return fmt.Errorf("%w: %w", errFailedToLstat, err)
	}

	if srcInfo.Mode()&os.ModeSymlink != 0 {
		linkTarget, err := os.Readlink(src)
		if err != nil {
			return fmt.Errorf("%w: %w", errFailedToWalkDir, err)
		}
		if err := os.Symlink(linkTarget, dst); err != nil {
			return fmt.Errorf("%w: %w", errFailedToCreateSymlink, err)
		}
		return nil
	}

	srcFile, err := os.Open(src) // #nosec G304 -- intended
	if err != nil {
		return fmt.Errorf("%w: %w", errFailedToOpenFile, err)
	}
	defer errorsx.Close(srcFile, &e, "src file")

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode().Perm()) // #nosec G304 -- intended
	if err != nil {
		return fmt.Errorf("%w: %w", errFailedToOpenFile, err)
	}
	defer errorsx.Close(dstFile, &e, "dst file")

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("%w: %w", errFailedToCopyFile, err)
	}

	return nil
}

func (parts archiveParts) restoreTriplestore(cmd *cobra.Command, instance *wisski.WissKI) (e error) {
	liquid := ingredient.GetLiquid(instance.TRB())

	stack, err := instance.Barrel().OpenStack()
	if err != nil {
		return fmt.Errorf("%w: %w", errFailedToOpenStack, err)
	}
	defer errorsx.Close(stack, &e, "stack")

	if _, err := logging.LogMessage(cmd.OutOrStdout(), "Purging triplestore repository"); err != nil {
		return fmt.Errorf("%w: %w", errFailedToLogMessage, err)
	}
	if err := liquid.TS.Purge(cmd.Context(), liquid.Instance, liquid.Domain()); err != nil {
		return fmt.Errorf("%w: %w", errFailedToPurgeTriplestoreData, err)
	}

	if _, err := logging.LogMessage(cmd.OutOrStdout(), "Re-provisioning triplestore repository"); err != nil {
		return fmt.Errorf("%w: %w", errFailedToLogMessage, err)
	}
	if err := liquid.TS.Provision(cmd.Context(), liquid.Instance, liquid.Domain(), &stack); err != nil {
		return fmt.Errorf("%w: %w", errFailedToProvisionTriplestore, err)
	}

	if _, err := logging.LogMessage(cmd.OutOrStdout(), "Restoring triplestore contents"); err != nil {
		return fmt.Errorf("%w: %w", errFailedToLogMessage, err)
	}

	file, err := os.Open(parts.TSFilePath)
	if err != nil {
		return fmt.Errorf("%w: %w", errFailedToOpenTriplestoreBackup, &fs.PathError{Op: "open", Path: parts.TSFilePath, Err: err})
	}
	defer func() {
		_ = file.Close()
	}()
	if err := liquid.TS.RestoreDB(cmd.Context(), liquid.GraphDBRepository, file); err != nil {
		return fmt.Errorf("%w: %w", errFailedToRestoreTriplestoreContents, err)
	}
	return nil
}

func (parts archiveParts) restoreSQL(cmd *cobra.Command, instance *wisski.WissKI) (e error) {
	liquid := ingredient.GetLiquid(instance.TRB())

	if _, err := logging.LogMessage(cmd.OutOrStdout(), "Purging SQL database"); err != nil {
		return fmt.Errorf("%w: %w", errFailedToLogMessage, err)
	}
	if err := liquid.BoundSQL().Purge(cmd.Context()); err != nil {
		return fmt.Errorf("%w: %w", errFailedToPurgeSQLDatabase, err)
	}

	if _, err := logging.LogMessage(cmd.OutOrStdout(), "Re-provisioning SQL database"); err != nil {
		return fmt.Errorf("%w: %w", errFailedToLogMessage, err)
	}
	if err := liquid.BoundSQL().Provision(cmd.Context()); err != nil {
		return fmt.Errorf("%w: %w", errFailedToProvisionSQLDatabase, err)
	}

	if _, err := logging.LogMessage(cmd.OutOrStdout(), "Restoring SQL contents"); err != nil {
		return fmt.Errorf("%w: %w", errFailedToLogMessage, err)
	}

	file, err := os.Open(parts.SQLFilePath)
	if err != nil {
		return fmt.Errorf("%w: %w", errFailedToOpenSQLBackup, &fs.PathError{Op: "open", Path: parts.SQLFilePath, Err: err})
	}
	defer errorsx.Close(file, &e, "file")

	if err := liquid.BoundSQL().Restore(cmd.Context(), file, stream.NewIOStream(cmd.OutOrStdout(), cmd.ErrOrStderr(), nil)); err != nil {
		return fmt.Errorf("%w: %w", errFailedToRestoreSQLContents, err)
	}
	return nil
}

func (parts archiveParts) restoreSQLConfig(cmd *cobra.Command, instance *wisski.WissKI) (e error) {
	if err := instance.Settings().SetDefaultDBConnection(cmd.Context(), nil, instance.BoundSQL().SQLUrl()); err != nil {
		return fmt.Errorf("%w: %w", errFailedToRestoreSQLConfig, err)
	}
	return nil
}
