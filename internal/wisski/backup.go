package wisski

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/countwriter"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/exp/slices"
)

// backupDescription is a description for a backup
type BackupDescription struct {
	Dest string // destination path
}

// Snapshot represents the result of generating a snapshot
type Backup struct {
	Description BackupDescription

	// Start and End Time of the backup
	StartTime time.Time
	EndTime   time.Time

	// various error states, which are ignored when creating the snapshot
	ErrPanic interface{}

	// errors for the various components
	ComponentErrors map[string]error

	// SQL and triplestore errors
	TSErr error

	// TODO: Make this proper
	ConfigFileErr error

	// Snapshots containing instances
	InstanceListErr   error
	InstanceSnapshots []Snapshot

	// List of files included
	Manifest []string
}

func (backup Backup) String() string {
	var builder strings.Builder
	backup.Report(&builder)
	return builder.String()
}

// Report writes a report from backup into w
func (backup Backup) Report(w io.Writer) (int, error) {
	cw := countwriter.NewCountWriter(w)

	encoder := json.NewEncoder(cw)
	encoder.SetIndent("", "  ")

	io.WriteString(cw, "======= Backup =======\n")

	fmt.Fprintf(cw, "Start: %s\n", backup.StartTime)
	fmt.Fprintf(cw, "End:   %s\n", backup.EndTime)
	io.WriteString(cw, "\n")

	io.WriteString(cw, "======= Description =======\n")
	encoder.Encode(backup.Description)
	io.WriteString(cw, "\n")

	io.WriteString(cw, "======= Errors =======\n")
	fmt.Fprintf(cw, "Panic:            %v\n", backup.ErrPanic)
	fmt.Fprintf(cw, "Component Errors: %v\n", backup.ComponentErrors)
	fmt.Fprintf(cw, "ConfigFileErr:    %s\n", backup.ConfigFileErr)
	fmt.Fprintf(cw, "InstanceListErr:  %s\n", backup.InstanceListErr)

	io.WriteString(cw, "\n")

	io.WriteString(cw, "======= Snapshots =======\n")
	for _, s := range backup.InstanceSnapshots {
		io.WriteString(cw, s.String())
		io.WriteString(cw, "\n")
	}

	io.WriteString(cw, "======= Manifest =======\n")
	for _, file := range backup.Manifest {
		io.WriteString(cw, file+"\n")
	}

	io.WriteString(cw, "\n")

	return cw.Sum()
}

// Backup makes a makes of the entire distillery.
// To make a backup, all [BackupComponents] will be invoked.
func (dis *Distillery) Backup(io stream.IOStream, description BackupDescription) (backup Backup) {
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

func (backup *Backup) run(io stream.IOStream, dis *Distillery) {

	backups := dis.Backupable()

	files := make(chan string, len(backups))         // channel for files being added into the backups
	results := make(chan backupResult, len(backups)) // channel for results to be stored into
	backup.ComponentErrors = make(map[string]error, len(backups))

	wg := &sync.WaitGroup{} // to wait for the results
	wg.Add(len(backups))
	for _, bc := range backups {
		go func(bc component.Backupable) {
			defer wg.Done()

			// find the backup destination
			dest := filepath.Join(backup.Description.Dest, bc.BackupName())
			files <- dest

			// make the backup and send the result!
			results <- backupResult{
				name: bc.Name(),
				err:  bc.Backup(io, dest),
			}
		}(bc)
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

		backup.InstanceSnapshots = make([]Snapshot, len(instances))
		for i, instance := range instances {
			backup.InstanceSnapshots[i] = func() Snapshot {
				dir := filepath.Join(instancesBackupDir, instance.Slug)
				if err := os.Mkdir(dir, fs.ModeDir); err != nil {
					return Snapshot{
						ErrPanic: err,
					}
				}

				files <- dir
				return dis.Snapshot(instance, io.NonInteractive(), SnapshotDescription{
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
	slices.SortFunc(backup.InstanceSnapshots, func(a, b Snapshot) bool {
		return a.Instance.Slug < b.Instance.Slug
	})
}

// WriteReport writes out the report belonging to this backup.
// It is a separate function, to allow writing it indepenently of the rest.
func (backup Backup) WriteReport(io stream.IOStream) error {
	return logging.LogOperation(func() error {
		reportPath := filepath.Join(backup.Description.Dest, "report.txt")
		io.Println(reportPath)

		// create the report file!
		report, err := os.Create(reportPath)
		if err != nil {
			return err
		}
		defer report.Close()

		// print the report into it!
		_, err = report.WriteString(backup.String())
		return err
	}, io, "Writing backup report")
}

// ShouldPrune determines if a file with the provided modtime
func (dis *Distillery) ShouldPrune(modtime time.Time) bool {
	return time.Since(modtime) > time.Duration(dis.Config.MaxBackupAge)*24*time.Hour
}

// PruneBackups prunes all backups older than the maximum backup age
func (dis *Distillery) PruneBackups(io stream.IOStream) error {
	sPath := dis.SnapshotsArchivePath()

	// list all the files
	entries, err := os.ReadDir(sPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// skip directories
		if entry.IsDir() {
			continue
		}

		// grab info about the file
		info, err := entry.Info()
		if err != nil {
			return err
		}

		// check if it should be pruned!
		if !dis.ShouldPrune(info.ModTime()) {
			continue
		}

		// assemble path, and then remove the file!
		path := filepath.Join(sPath, entry.Name())
		io.Printf("Removing %s cause it is older than %d days", path, dis.Config.MaxBackupAge)

		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}
