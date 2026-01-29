package cmd

//spellchecker:words bufio encoding json path filepath strings github wisski distillery internal component exporter cobra pkglib exit
import (
	"bufio"
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
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

	checkResult, err := readArchiveParts(cmd, sr.Positionals.Directory, snapshot)
	if err != nil {
		return fmt.Errorf("snapshot is not suitable for restoration: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "================================================\n")
	fmt.Fprintf(cmd.OutOrStdout(), "Instance: %s\n", instance.FilesystemBase)
	fmt.Fprintf(cmd.OutOrStdout(), "Snapshot: %s\n", sr.Positionals.Directory)
	fmt.Fprintf(cmd.OutOrStdout(), "================================================\n")
	fmt.Fprintf(cmd.OutOrStdout(), "Data:        %s\n", checkResult.DataPath)
	fmt.Fprintf(cmd.OutOrStdout(), "SQL:         %s\n", checkResult.SQLFilePath)
	fmt.Fprintf(cmd.OutOrStdout(), "Triplestore: %s\n", checkResult.TSFilePath)
	fmt.Fprintf(cmd.OutOrStdout(), "================================================\n")

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

		fmt.Fprintf(cmd.OutOrStdout(), "Old data directory: %s\n", oldDataDirectory)
		fmt.Fprintf(cmd.OutOrStdout(), "New data directory: %s\n", checkResult.DataPath)

		if err := checkResult.restoreDataDirectory(cmd, oldDataDirectory); err != nil {
			return fmt.Errorf("failed to restore data directory: %w", err)
		}

		// TODO: chown the data directory to the www-data user

		fmt.Fprintln(cmd.OutOrStdout(), "Data directory restored.")
	}

	// Triplestore
	{
		if err := logging.LogOperation(func() error {
			return checkResult.restoreTriplestore(cmd, instance)
		}, cmd.ErrOrStderr(), "Restoring triplestore"); err != nil {
			return fmt.Errorf("failed to restore triplestore: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Triplestore restored.")
	}

	// SQL
	{
		if err := logging.LogOperation(func() error {
			return checkResult.restoreSQL(cmd, dis, instance)
		}, cmd.ErrOrStderr(), "Restoring SQL"); err != nil {
			return fmt.Errorf("failed to restore SQL database: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "SQL restored.")
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
		fmt.Fprintln(cmd.OutOrStdout(), "Re-setting permissions and ownership")

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

		if err := checkResult.restoreSQLConfig(cmd, dis, instance); err != nil {
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
		fmt.Fprintln(cmd.OutOrStdout(), "Re-Creating Adapter")

		if _, err := instance.Adapters().SetAdapter(cmd.Context(), nil, instance.Adapters().DefaultAdapter()); err != nil {
			return fmt.Errorf("failed to restore adapter: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Adapter re-created.")
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

	// clear cache
	{
		fmt.Fprintln(cmd.OutOrStdout(), "Clearing cache.")

		if err := instance.Drush().Exec(cmd.Context(), cmd.OutOrStdout(), "cr"); err != nil {
			return fmt.Errorf("failed to run drush cr: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Cache cleared.")
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

	fmt.Fprintln(cmd.OutOrStdout(), "Instance should be restored.")
	return nil
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

type archiveParts struct {
	DataPath string

	SQLFilePath     string
	SQLDatabaseName string

	TSFilePath string
}

func readArchiveParts(cmd *cobra.Command, path string, archive exporter.Snapshot) (parts archiveParts, err error) {
	{
		if !slices.Contains(archive.Description.Parts, "data") {
			return archiveParts{}, fmt.Errorf("data part is not present in the snapshot")
		}

		parts.DataPath = filepath.Join(path, "data", "data")
		if isDirectory, err := fsx.IsDirectory(parts.DataPath, false); !isDirectory {
			return archiveParts{}, fmt.Errorf("data folder is not a directory: %w", cmp.Or(err, fs.ErrNotExist))
		}
	}

	{
		if !slices.Contains(archive.Description.Parts, "sql") {
			return archiveParts{}, fmt.Errorf("sql part missing from snapshot manifest")
		}

		parts.SQLFilePath = filepath.Join(path, "sql", archive.Instance.SqlDatabase+".sql")
		if isFile, err := fsx.IsRegular(parts.SQLFilePath, false); !isFile {
			return archiveParts{}, fmt.Errorf("sql data not found in snapshot: %w", cmp.Or(err, fs.ErrNotExist))
		}
		parts.SQLDatabaseName = archive.Instance.SqlDatabase
	}

	{
		if !slices.Contains(archive.Description.Parts, "triplestore") {
			return archiveParts{}, fmt.Errorf("triplestore part missing from snapshot manifest")
		}

		parts.TSFilePath = filepath.Join(path, "triplestore", archive.Instance.GraphDBRepository+".nq")
		if isFile, err := fsx.IsRegular(parts.TSFilePath, false); !isFile {
			return archiveParts{}, fmt.Errorf("triplestore data not a regular file: %w", cmp.Or(err, fs.ErrNotExist))
		}
	}

	return parts, nil
}

func shutdownInstance(cmd *cobra.Command, instance *wisski.WissKI) (e error) {
	stack, err := instance.Barrel().OpenStack()
	if err != nil {
		return fmt.Errorf("failed to open stack: %w", err)
	}
	defer errorsx.Close(stack, &e, "stack")
	return stack.Down(cmd.Context(), cmd.OutOrStdout())
}

func startInstance(cmd *cobra.Command, instance *wisski.WissKI) (e error) {
	stack, err := instance.Barrel().OpenStack()
	if err != nil {
		return fmt.Errorf("failed to open stack: %w", err)
	}
	defer errorsx.Close(stack, &e, "stack")
	return stack.Start(cmd.Context(), cmd.OutOrStdout())
}

func (parts archiveParts) restoreDataDirectory(cmd *cobra.Command, old string) (e error) {
	new := parts.DataPath

	st := status.NewWithCompat(cmd.ErrOrStderr(), 1)
	st.Start()
	defer st.Set(0, "")
	defer st.Stop()

	// Create a temporary directory next to the old directory
	oldParent := filepath.Dir(old)
	tempDir, err := os.MkdirTemp(oldParent, ".restore-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}

	// Clean up temporary directory on failure
	defer func() {
		if e != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "failed to restore, cleaning up")
			os.RemoveAll(tempDir)
		}
	}()

	// Copy the new directory contents to the temporary directory
	if err := filepath.WalkDir(new, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Compute relative path and destination
		relPath, err := filepath.Rel(new, path)
		if err != nil {
			return fmt.Errorf("failed to compute relative path: %w", err)
		}
		destPath := filepath.Join(tempDir, relPath)

		st.Set(0, relPath)

		if d.IsDir() {
			// Create directory with same permissions
			info, err := d.Info()
			if err != nil {
				return fmt.Errorf("failed to get directory info: %w", err)
			}
			if err := os.MkdirAll(destPath, info.Mode().Perm()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
			return nil
		}

		// Copy regular file
		if err := copyFile(path, destPath); err != nil {
			return fmt.Errorf("failed to copy file %s: %w", relPath, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to copy directory: %w", err)
	}

	// Remove the old directory
	if err := os.RemoveAll(old); err != nil {
		return fmt.Errorf("failed to remove old directory: %w", err)
	}

	// Move the temporary directory to the old directory's place
	if err := os.Rename(tempDir, old); err != nil {
		return fmt.Errorf("failed to move restored directory into place: %w", err)
	}

	return nil
}

// Copies a file from src to dst.
// If it is a symlink, it is copied as a symlink.
func copyFile(src, dst string) error {
	srcInfo, err := os.Lstat(src)
	if err != nil {
		return err
	}

	if srcInfo.Mode()&os.ModeSymlink != 0 {
		linkTarget, err := os.Readlink(src)
		if err != nil {
			return err
		}
		return os.Symlink(linkTarget, dst)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode().Perm())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}

func (parts archiveParts) restoreTriplestore(cmd *cobra.Command, instance *wisski.WissKI) (e error) {
	liquid := ingredient.GetLiquid(instance.TRB())

	if _, err := logging.LogMessage(cmd.OutOrStdout(), "Purging triplestore repository"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if err := liquid.TS.Purge(cmd.Context(), liquid.Instance, liquid.Domain()); err != nil {
		return fmt.Errorf("failed to purge triplestore data: %w", err)
	}

	if _, err := logging.LogMessage(cmd.OutOrStdout(), "Re-provisioning triplestore repository"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if err := liquid.TS.Provision(cmd.Context(), liquid.Instance, liquid.Domain()); err != nil {
		return fmt.Errorf("failed to provision triplestore: %w", err)
	}

	if _, err := logging.LogMessage(cmd.OutOrStdout(), "Restoring triplestore contents"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}

	file, err := os.Open(parts.TSFilePath)
	if err != nil {
		return fmt.Errorf("failed to open triplestore backup: %w", err)
	}
	defer file.Close()

	if err := liquid.TS.RestoreDB(cmd.Context(), liquid.GraphDBRepository, file); err != nil {
		return fmt.Errorf("failed to restore triplestore contents: %w", err)
	}
	return nil
}

func (parts archiveParts) restoreSQL(cmd *cobra.Command, dis *dis.Distillery, instance *wisski.WissKI) (e error) {
	liquid := ingredient.GetLiquid(instance.TRB())

	if _, err := logging.LogMessage(cmd.OutOrStdout(), "Purging SQL database"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if err := liquid.SQL.Purge(cmd.Context(), liquid.Instance, liquid.Domain()); err != nil {
		return fmt.Errorf("failed to purge SQL database: %w", err)
	}

	if _, err := logging.LogMessage(cmd.OutOrStdout(), "Re-provisioning SQL database"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if err := liquid.SQL.Provision(cmd.Context(), liquid.Instance, liquid.Domain()); err != nil {
		return fmt.Errorf("failed to provision SQL database: %w", err)
	}

	if _, err := logging.LogMessage(cmd.OutOrStdout(), "Restoring SQL contents"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}

	file, err := os.Open(parts.SQLFilePath)
	if err != nil {
		return fmt.Errorf("failed to open SQL backup: %w", err)
	}
	defer file.Close()

	replacedFile := replaceSqlDatabaseName(file, instance.SqlDatabase, parts.SQLDatabaseName)
	defer replacedFile.Close()

	//
	code := dis.SQL().Shell(cmd.Context(), stream.NewIOStream(cmd.OutOrStdout(), cmd.ErrOrStderr(), replacedFile))
	if code != 0 {
		return fmt.Errorf("failed to restore SQL contents: exit code %d", code)
	}
	return nil
}

var (
	reCreateDB = regexp.MustCompile(
		`(?i)^\s*CREATE\s+DATABASE\b.*?` + "`" + `([^` + "`" + `]+)` + "`",
	)
	reUseDB = regexp.MustCompile(
		`(?i)^\s*USE\s+(` + "`" + `([^` + "`" + `]+)` + "`" + `|([^\s;]+))`,
	)
)

func replaceSqlDatabaseName(reader io.Reader, newDB string, oldDB string) io.ReadCloser {
	// HACK HACK HACK: This restore code makes shit tons of assumptions about the SQL dump.
	// In particular that it was created by mysqldump -- and only one 'CREATE DATABASE' statement exists.

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		scanner := bufio.NewScanner(reader)
		// allow large lines (mysqldump can emit big INSERTs)
		scanner.Buffer(make([]byte, 64*1024), 10*1024*1024)

		for scanner.Scan() {
			line := scanner.Text()

			// CREATE DATABASE ... `oldDB`
			if m := reCreateDB.FindStringSubmatchIndex(line); m != nil {
				// group 1 = db name inside backticks
				dbStart, dbEnd := m[2], m[3]
				if line[dbStart:dbEnd] == oldDB {
					line = line[:dbStart] + newDB + line[dbEnd:]
				}
			} else if m := reUseDB.FindStringSubmatchIndex(line); m != nil {
				// USE `db`  -> group 2
				// USE db    -> group 3
				dbStart, dbEnd := -1, -1
				if m[4] != -1 { // backticked
					dbStart, dbEnd = m[4], m[5]
				} else if m[6] != -1 { // bare
					dbStart, dbEnd = m[6], m[7]
				}
				if dbStart != -1 && line[dbStart:dbEnd] == oldDB {
					line = line[:dbStart] + newDB + line[dbEnd:]
				}
			}

			if _, err := io.WriteString(pw, line+"\n"); err != nil {
				_ = pw.CloseWithError(err)
				return
			}
		}

		if err := scanner.Err(); err != nil {
			_ = pw.CloseWithError(err)
		}
	}()

	return pr
}

func (parts archiveParts) restoreSQLConfig(cmd *cobra.Command, dis *dis.Distillery, instance *wisski.WissKI) (e error) {
	if err := instance.Settings().SetDefaultDBConnection(cmd.Context(), nil, instance.SQLURL()); err != nil {
		return fmt.Errorf("failed to restore SQL config: %w", err)
	}
	return nil
}
