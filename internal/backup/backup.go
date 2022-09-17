// Package backup implements Distillery backups.
package backup

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/exp/slices"
)

// New create a new Backup
func New(io stream.IOStream, dis *wisski.Distillery, description Description) (backup Backup) {
	backup.Description = description

	// catch anything critical that happened during the snapshot
	defer func() {
		backup.ErrPanic = recover()
	}()

	// do the create keeping track of time!
	logging.LogOperation(func() error {
		backup.StartTime = time.Now()
		backup.run(io, dis)
		backup.EndTime = time.Now()

		return nil
	}, io, "Writing backup files")

	return
}

type backupResult struct {
	name string
	err  error
}

func (backup *Backup) run(io stream.IOStream, dis *wisski.Distillery) {

	backups := dis.Backupable()

	files := make(chan string, len(backups))         // channel for files being added into the backups
	results := make(chan backupResult, len(backups)) // channel for results to be stored into
	backup.ComponentErrors = make(map[string]error, len(backups))

	wg := &sync.WaitGroup{} // to wait for the results
	wg.Add(len(backups))    // tell the group about all the operations
	for _, bc := range backups {
		go func(bc component.Backupable, context component.BackupContext) {
			defer wg.Done()

			// make the backup and store the result
			results <- backupResult{
				name: bc.Name(),
				err:  bc.Backup(context),
			}
		}(bc, &context{
			io:    io,
			dst:   filepath.Join(backup.Description.Dest, bc.BackupName()),
			files: files,
		})
	}

	// backup instances
	wg.Add(1)
	go func() {
		defer wg.Done()

		instancesBackupDir := filepath.Join(backup.Description.Dest, "instances")
		if err := os.Mkdir(instancesBackupDir, fs.ModeDir); err != nil {
			backup.InstanceListErr = err
			return
		}

		// list all instances
		instances, err := dis.Instances().All()
		if err != nil {
			backup.InstanceListErr = err
			return
		}

		backup.InstanceSnapshots = make([]wisski.Snapshot, len(instances))
		for i, instance := range instances {
			backup.InstanceSnapshots[i] = func() wisski.Snapshot {
				dir := filepath.Join(instancesBackupDir, instance.Slug)
				if err := os.Mkdir(dir, fs.ModeDir); err != nil {
					return wisski.Snapshot{
						ErrPanic: err,
					}
				}

				files <- dir
				return dis.Snapshot(instance, io.NonInteractive(), wisski.SnapshotDescription{
					Dest: dir,
				})
			}()
		}

	}()

	// finish processing all the results as soon as the group is done.
	go func() {
		defer close(results)
		wg.Wait()
	}()

	// finish the message processing once results are finished.
	go func() {
		defer close(files) // no more file processing!
		for result := range results {
			backup.ComponentErrors[result.name] = result.err
		}
	}()

	for file := range files {
		// get the relative path to the root of the manifest.
		// nothing *should* go wrong, but in case it does, use the original path.
		path, err := filepath.Rel(backup.Description.Dest, file)
		if err != nil {
			path = file
		}

		// write it to the command line
		// and also add it to the manifest
		io.Printf("\033[2K\r%s", path)
		backup.Manifest = append(backup.Manifest, path)
	}
	slices.Sort(backup.Manifest) // backup the manifest
	io.Println("")

	// sort the instances manifest
	slices.SortFunc(backup.InstanceSnapshots, func(a, b wisski.Snapshot) bool {
		return a.Instance.Slug < b.Instance.Slug
	})
}
