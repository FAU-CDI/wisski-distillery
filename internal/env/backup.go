package env

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/pkg/errors"
	"github.com/tkw1536/goprogram/stream"
)

// backupDescription is a description for a backup
type BackupDescription struct {
	Dest string // destination path
}

// Snapshot represents the result of generating a snapshot
type Backup struct {
	Description BackupDescription

	// various error states, which are ignored when creating the snapshot
	ErrPanic interface{}

	SQLErr error
	TSErr  error

	ConfigFileErr       error
	ConfigFilesManifest map[string]error

	InstanceListErr   error
	InstancesManifest []Snapshot
}

func (dis *Distillery) Backup(io stream.IOStream, description BackupDescription) (backup Backup) {
	backup.Description = description

	// catch anything critical that happened during the snapshot
	defer func() {
		backup.ErrPanic = recover()
	}()

	backup.run(io, dis)
	return
}

var errBackupSkipFile = errors.New("<file not found>")

func (backup *Backup) run(io stream.IOStream, dis *Distillery) {
	// create a wait group, and message channel
	wg := &sync.WaitGroup{}
	messages := make(chan string, 4)

	// backup the sql
	wg.Add(1)
	go func() {
		defer wg.Done()

		sqlPath := filepath.Join(backup.Description.Dest, "sql.sql")
		messages <- sqlPath

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
		messages <- tsPath

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

		files := []string{
			filepath.Join(dis.Config.DeployRoot, ".env"),  // TODO: put the name of the configuration file somewhere
			filepath.Join(dis.Config.DeployRoot, "wdcli"), // TODO: constant the name of the executable
			dis.Config.SelfOverridesFile,
			dis.Config.GlobalAuthorizedKeysFile,
		}

		backup.ConfigFilesManifest = make(map[string]error, len(files))
		for _, src := range files {
			if !fsx.IsFile(src) {
				backup.ConfigFilesManifest[src] = errBackupSkipFile
				continue
			}
			dest := filepath.Join(cfgBackupDir, filepath.Base(src))

			// copy the config file and store the error message
			messages <- src
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
		instances, err := dis.AllInstances()
		if err != nil {
			backup.InstanceListErr = err
			return
		}

		iochild := stream.NewIOStream(io.Stderr, io.Stderr, nil, 0)

		backup.InstancesManifest = make([]Snapshot, len(instances))
		for i, instance := range instances {
			backup.InstancesManifest[i] = func() Snapshot {
				dir := filepath.Join(instancesBackupDir, instance.Slug)
				if err := os.Mkdir(dir, fs.ModeDir); err != nil {
					return Snapshot{
						ErrPanic: err,
					}
				}

				messages <- dir
				return instance.Snapshot(iochild, SnapshotDescription{
					Dest: dir,
				})
			}()
		}

	}()

	// wait for the group, then close the message channel.
	go func() {
		wg.Wait()
		close(messages)
	}()

	// print out all the messages
	for message := range messages {
		io.Println(message)
	}
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
		_, err = fmt.Fprintf(report, "%#v\n", backup)
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
