//spellchecker:words exporter
package exporter

//spellchecker:words context errors slog path filepath strings time github wisski distillery internal component wdlog logging pkglib umaskfree status golang slices
import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"

	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/pkglib/errorsx"
	"github.com/tkw1536/pkglib/fsx/umaskfree"
	"github.com/tkw1536/pkglib/status"
)

// Backup describes a backup.
//
//nolint:recvcheck
type Backup struct {
	Description BackupDescription

	// Start and End Time of the backup
	StartTime time.Time
	EndTime   time.Time

	// various error states, which are ignored when creating the snapshot
	ErrPanic interface{}

	// errors for the various components
	ComponentErrors map[string]error

	// TODO: Make this proper
	ConfigFileErr error

	// Snapshots containing instances
	InstanceListErr   error
	InstanceSnapshots []Snapshot

	// List of files included
	WithManifest
}

// BackupDescription provides a description for a backup.
type BackupDescription struct {
	Dest string // Destination path

	ConcurrentSnapshots int // maximum number of concurrent snapshots
}

// New create a new Backup.
func (exporter *Exporter) NewBackup(ctx context.Context, progress io.Writer, description BackupDescription) (backup Backup) {
	backup.Description = description

	// catch anything critical that happened during the snapshot
	defer func() {
		backup.ErrPanic = recover()
	}()

	// do the create keeping track of time!
	err := logging.LogOperation(func() error {
		backup.StartTime = time.Now().UTC()
		err := backup.run(ctx, progress, exporter)
		backup.EndTime = time.Now().UTC()

		return err
	}, progress, "Writing backup files")

	if err != nil {
		wdlog.Of(ctx).Error("failed to create new backup (using possibly incomplete result)", slog.Any("error", err))
	}

	return backup
}

func (backup *Backup) run(ctx context.Context, progress io.Writer, exporter *Exporter) error {
	// create a manifest
	manifest, done := backup.handleManifest(backup.Description.Dest)
	defer done()

	// create a new status display
	backups := exporter.dependencies.Backupable
	backup.ComponentErrors = make(map[string]error, len(backups))

	// Component backup tasks
	err1 := logging.LogOperation(func() error {
		st := status.NewWithCompat(progress, 0)
		st.Start()
		defer st.Stop()

		errors, _ := status.Group[component.Backupable, error]{
			PrefixString: func(item component.Backupable, index int) string {
				return fmt.Sprintf("[backup %q]: ", item.Name())
			},
			PrefixAlign: true,

			Handler: func(bc component.Backupable, index int, writer io.Writer) error {
				return bc.Backup(
					component.NewStagingContext(
						ctx,
						writer,
						filepath.Join(backup.Description.Dest, bc.BackupName()),
						manifest,
					),
				)
			},

			ResultString: status.DefaultErrorString[component.Backupable],
		}.Use(st, backups)

		for i, bc := range backups {
			backup.ComponentErrors[bc.Name()] = errors[i]
		}

		return nil
	}, progress, "Backing up core components")

	// backup instances
	err2 := logging.LogOperation(func() error {
		st := status.NewWithCompat(progress, 0)
		st.Start()
		defer st.Stop()

		instancesBackupDir := filepath.Join(backup.Description.Dest, "instances")
		if err := umaskfree.Mkdir(instancesBackupDir, umaskfree.DefaultDirPerm); err != nil {
			backup.InstanceListErr = err
			return nil
		}

		// list all instances
		wissKIs, err := exporter.dependencies.Instances.All(ctx)
		if err != nil {
			backup.InstanceListErr = err
			return nil
		}

		// make a backup of the snapshots
		backup.InstanceSnapshots, _ = status.Group[*wisski.WissKI, Snapshot]{
			PrefixString: func(item *wisski.WissKI, index int) string {
				return fmt.Sprintf("[snapshot %q]: ", item.Slug)
			},
			PrefixAlign: true,

			Handler: func(instance *wisski.WissKI, index int, writer io.Writer) Snapshot {
				dir := filepath.Join(instancesBackupDir, instance.Slug)
				if err := umaskfree.Mkdir(dir, umaskfree.DefaultDirPerm); err != nil {
					return Snapshot{
						ErrPanic: err,
					}
				}

				manifest <- dir

				return exporter.NewSnapshot(ctx, instance, writer, SnapshotDescription{
					Dest: dir,
				})
			},
			ResultString: func(res Snapshot, item *wisski.WissKI, index int) string {
				return "done"
			},
			WaitString:   status.DefaultWaitString[*wisski.WissKI],
			HandlerLimit: backup.Description.ConcurrentSnapshots,
		}.Use(st, wissKIs)

		// sort the instances
		slices.SortFunc(backup.InstanceSnapshots, func(a, b Snapshot) int {
			return strings.Compare(a.Instance.Slug, b.Instance.Slug)
		})

		return nil
	}, progress, "Creating instance snapshots")

	return errorsx.Combine(err1, err2)
}
