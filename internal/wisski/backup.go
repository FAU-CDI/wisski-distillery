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

	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/pkg/countwriter"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/pkg/errors"
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

	// SQL and triplestore errors
	SQLErr error
	TSErr  error

	// TODO: Make this proper
	ConfigFileErr       error
	ConfigFilesManifest map[string]error

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
	fmt.Fprintf(cw, "Panic:           %v\n", backup.ErrPanic)
	fmt.Fprintf(cw, "SQLErr:          %s\n", backup.SQLErr)
	fmt.Fprintf(cw, "TSErr:           %s\n", backup.TSErr)
	fmt.Fprintf(cw, "ConfigFileErr:   %s\n", backup.ConfigFileErr)
	fmt.Fprintf(cw, "InstanceListErr: %s\n", backup.InstanceListErr)

	io.WriteString(cw, "\n")

	io.WriteString(cw, "======= Config Files =======\n")
	encoder.Encode(backup.ConfigFilesManifest) // TODO: Proper manifest

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

var errBackupSkipFile = errors.New("<file not found>")

func (backup *Backup) run(io stream.IOStream, dis *Distillery) {
	// create a wait group, and message channel
	wg := &sync.WaitGroup{}
	files := make(chan string, 4)

	// backup the sql
	wg.Add(1)
	go func() {
		defer wg.Done()

		sqlPath := filepath.Join(backup.Description.Dest, "sql.sql")
		files <- sqlPath

		sql, err := os.Create(sqlPath)
		if err != nil {
			backup.SQLErr = err
			return
		}
		defer sql.Close()

		// directly store the result
		backup.SQLErr = dis.SQL().BackupAll(io, sql)
	}()

	// backup the triplestore
	wg.Add(1)
	go func() {
		defer wg.Done()

		tsPath := filepath.Join(backup.Description.Dest, "triplestore")
		files <- tsPath

		// directly store the result
		backup.TSErr = dis.Triplestore().BackupAll(tsPath)
	}()

	// backup configuration files
	wg.Add(1)
	go func() {
		defer wg.Done()

		cfgBackupDir := filepath.Join(backup.Description.Dest, "config")
		if err := os.Mkdir(cfgBackupDir, fs.ModeDir); err != nil {
			backup.ConfigFileErr = err
			return
		}

		configs := []string{
			dis.Config.ConfigPath,
			filepath.Join(dis.Config.DeployRoot, core.Executable), // TODO: constant the name of the executable
			dis.Config.SelfOverridesFile,
			dis.Config.GlobalAuthorizedKeysFile,
		}

		backup.ConfigFilesManifest = make(map[string]error, len(configs))
		for _, src := range configs {
			if !fsx.IsFile(src) {
				backup.ConfigFilesManifest[src] = errBackupSkipFile
				continue
			}
			dest := filepath.Join(cfgBackupDir, filepath.Base(src))

			// copy the config file and store the error message
			files <- src
			backup.ConfigFilesManifest[src] = fsx.CopyFile(dest, src)
		}
	}()

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

	// wait for the group, then close the message channel.
	go func() {
		wg.Wait()
		close(files)
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
