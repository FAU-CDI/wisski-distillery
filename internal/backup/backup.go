// Package backup implements Distillery backups.
package backup

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/component/snapshots"
	"github.com/FAU-CDI/wisski-distillery/internal/dis"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/status"
	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/exp/slices"
)

// New create a new Backup
func New(io stream.IOStream, dis *dis.Distillery, description Description) (backup Backup) {
	backup.Description = description

	// catch anything critical that happened during the snapshot
	defer func() {
		backup.ErrPanic = recover()
	}()

	// do the create keeping track of time!
	logging.LogOperation(func() error {
		backup.StartTime = time.Now().UTC()
		backup.run(io, dis)
		backup.EndTime = time.Now().UTC()

		return nil
	}, io, "Writing backup files")

	return
}

func (backup *Backup) run(ios stream.IOStream, dis *dis.Distillery) {
	//
	// MANIFEST
	//

	manifest := make(chan string)       // receive all the entries in the manifest
	manifestDone := make(chan struct{}) // to signal that everything is finished

	go func() {
		defer close(manifestDone)

		for file := range manifest {
			// get the relative path to the root of the manifest
			// or fallback to the absolute path!
			path, err := filepath.Rel(backup.Description.Dest, file)
			if err != nil {
				path = file
			}

			// add the file to the manifest array
			backup.Manifest = append(backup.Manifest, path)
		}

		// sort the manifest
		slices.Sort(backup.Manifest)
	}()

	//
	// BACKUP COMPONENTS
	//

	// create a new status display
	backups := dis.Backupable()
	backup.ComponentErrors = make(map[string]error, len(backups))

	// Component backup tasks
	logging.LogOperation(func() error {
		st := status.NewWithCompat(ios.Stdout, 0)
		st.Start()
		defer st.Stop()

		return status.UseErrorGroup(st, status.Group[component.Backupable, error]{
			PrefixString: func(item component.Backupable, index int) string {
				return fmt.Sprintf("[backup %q]: ", item.Name())
			},
			PrefixAlign: true,

			Handler: func(bc component.Backupable, index int, writer io.Writer) error {
				// create a new context for the backup!
				context := &context{
					env:   dis.Core.Environment,
					io:    stream.NewIOStream(writer, writer, nil, 0),
					dst:   filepath.Join(backup.Description.Dest, bc.BackupName()),
					files: manifest,
				}

				backup.ComponentErrors[bc.Name()] = bc.Backup(context)
				return nil
			},
		}, backups)
	}, ios, "Backing up core components")

	// backup instances
	logging.LogOperation(func() error {
		st := status.NewWithCompat(ios.Stdout, 0)
		st.Start()
		defer st.Stop()

		instancesBackupDir := filepath.Join(backup.Description.Dest, "instances")
		if err := dis.Core.Environment.Mkdir(instancesBackupDir, environment.DefaultDirPerm); err != nil {
			backup.InstanceListErr = err
			return nil
		}

		// list all instances
		wissKIs, err := dis.Instances().All()
		if err != nil {
			backup.InstanceListErr = err
			return nil
		}

		// re-use the backup of the snapshots
		backup.InstanceSnapshots = status.Group[instances.WissKI, snapshots.Snapshot]{
			PrefixString: func(item instances.WissKI, index int) string {
				return fmt.Sprintf("[snapshot %s]: ", item.Slug)
			},
			PrefixAlign: true,

			Handler: func(instance instances.WissKI, index int, writer io.Writer) snapshots.Snapshot {
				dir := filepath.Join(instancesBackupDir, instance.Slug)
				if err := dis.Core.Environment.Mkdir(dir, environment.DefaultDirPerm); err != nil {
					return snapshots.Snapshot{
						ErrPanic: err,
					}
				}

				manifest <- dir

				return dis.SnapshotManager().NewSnapshot(instance, stream.NewIOStream(writer, writer, nil, 0), snapshots.SnapshotDescription{
					Dest: dir,
				})
			},
			ResultString: func(res snapshots.Snapshot, item instances.WissKI, index int) string {
				return "done"
			},
			WaitString:   status.DefaultWaitString[instances.WissKI],
			HandlerLimit: backup.Description.ConcurrentSnapshots,
		}.Use(st, wissKIs)
		return nil
	}, ios, "creating instance snapshots")

	// close the manifest
	close(manifest)
	<-manifestDone

	// sort the instances manifest
	slices.SortFunc(backup.InstanceSnapshots, func(a, b snapshots.Snapshot) bool {
		return a.Instance.Slug < b.Instance.Slug
	})
}
