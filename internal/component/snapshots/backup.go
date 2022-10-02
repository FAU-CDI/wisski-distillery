package snapshots

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"

	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/status"
	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/exp/slices"
)

// Backup describes a backup
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

// BackupDescription provides a description for a backup
type BackupDescription struct {
	Dest string // Destination path
	Auto bool   // Was the path created automatically?

	ConcurrentSnapshots int // maximum number of concurrent snapshots
}

// New create a new Backup
func (manager *Manager) NewBackup(io stream.IOStream, description BackupDescription) (backup Backup) {
	backup.Description = description

	// catch anything critical that happened during the snapshot
	defer func() {
		backup.ErrPanic = recover()
	}()

	// do the create keeping track of time!
	logging.LogOperation(func() error {
		backup.StartTime = time.Now().UTC()
		backup.run(io, manager)
		backup.EndTime = time.Now().UTC()

		return nil
	}, io, "Writing backup files")

	return
}

func (backup *Backup) run(ios stream.IOStream, manager *Manager) {
	// create a manifest
	manifest, done := backup.handleManifest(backup.Description.Dest)
	defer done()

	// create a new status display
	backups := manager.Backupable
	backup.ComponentErrors = make(map[string]error, len(backups))

	// Component backup tasks
	logging.LogOperation(func() error {
		st := status.NewWithCompat(ios.Stdout, 0)
		st.Start()
		defer st.Stop()

		errors := status.Group[component.Backupable, error]{
			PrefixString: func(item component.Backupable, index int) string {
				return fmt.Sprintf("[backup %q]: ", item.Name())
			},
			PrefixAlign: true,

			Handler: func(bc component.Backupable, index int, writer io.Writer) error {
				return bc.Backup(
					component.NewStagingContext(
						manager.Environment,
						stream.NewIOStream(writer, writer, nil, 0),
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
	}, ios, "Backing up core components")

	// backup instances
	logging.LogOperation(func() error {
		st := status.NewWithCompat(ios.Stdout, 0)
		st.Start()
		defer st.Stop()

		instancesBackupDir := filepath.Join(backup.Description.Dest, "instances")
		if err := manager.Environment.Mkdir(instancesBackupDir, environment.DefaultDirPerm); err != nil {
			backup.InstanceListErr = err
			return nil
		}

		// list all instances
		wissKIs, err := manager.Instances.All()
		if err != nil {
			backup.InstanceListErr = err
			return nil
		}

		// make a backup of the snapshots
		backup.InstanceSnapshots = status.Group[instances.WissKI, Snapshot]{
			PrefixString: func(item instances.WissKI, index int) string {
				return fmt.Sprintf("[snapshot %q]: ", item.Slug)
			},
			PrefixAlign: true,

			Handler: func(instance instances.WissKI, index int, writer io.Writer) Snapshot {
				dir := filepath.Join(instancesBackupDir, instance.Slug)
				if err := manager.Environment.Mkdir(dir, environment.DefaultDirPerm); err != nil {
					return Snapshot{
						ErrPanic: err,
					}
				}

				manifest <- dir

				return manager.NewSnapshot(instance, stream.NewIOStream(writer, writer, nil, 0), SnapshotDescription{
					Dest: dir,
				})
			},
			ResultString: func(res Snapshot, item instances.WissKI, index int) string {
				return "done"
			},
			WaitString:   status.DefaultWaitString[instances.WissKI],
			HandlerLimit: backup.Description.ConcurrentSnapshots,
		}.Use(st, wissKIs)

		// sort the instances
		slices.SortFunc(backup.InstanceSnapshots, func(a, b Snapshot) bool {
			return a.Instance.Slug < b.Instance.Slug
		})

		return nil
	}, ios, "Creating instance snapshots")

}
