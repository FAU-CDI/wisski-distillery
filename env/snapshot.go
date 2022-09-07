package env

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/fsx"
	"github.com/FAU-CDI/wisski-distillery/internal/logging"
	"github.com/FAU-CDI/wisski-distillery/internal/password"
	"github.com/tkw1536/goprogram/stream"
)

// SnapshotsDir returns the path that contains all snapshot related data.
func (dis Distillery) SnapshotsDir() string {
	return filepath.Join(dis.Config.DeployRoot, "snapshots")
}

// SnapshotsStagingPath returns the path to the directory containing a temporary staging area for snapshots.
// Use NewSnapshotStagingDir to generate a new staging area.
func (dis Distillery) SnapshotsStagingPath() string {
	return filepath.Join(dis.SnapshotsDir(), "staging")
}

// SnapshotsArchivePath returns the path to the directory containing all exported archives.
// Use NewSnapshotArchivePath to generate a path to a new archive in this directory.
func (dis Distillery) SnapshotsArchivePath() string {
	return filepath.Join(dis.SnapshotsDir(), "archives")
}

// NewSnapshotArchivePath returns the path to a new archive with the provided prefix.
// The path is guaranteed to not exist.
func (dis Distillery) NewSnapshotArchivePath(prefix string) (path string) {
	// TODO: Consider moving these into a subdirectory with the provided prefix.
	for path == "" || fsx.Exists(path) {
		name := dis.newSnapshotName(prefix) + ".tar.gz"
		path = filepath.Join(dis.SnapshotsArchivePath(), name)
	}
	return
}

// newSnapshot name returns a new basename for a snapshot with the provided prefix.
// The name is guaranteed to be unique within this process.
func (Distillery) newSnapshotName(prefix string) string {
	suffix, _ := password.Password(64) // silently ignore any errors!
	return fmt.Sprintf("%s-%d-%s", prefix, time.Now().Unix(), suffix)
}

// NewSnapshotStagingDir returns the path to a new snapshot directory.
// The directory is guaranteed to have been freshly created.
func (dis Distillery) NewSnapshotStagingDir(prefix string) (path string, err error) {
	for path == "" || os.IsExist(err) {
		path = filepath.Join(dis.SnapshotsStagingPath(), dis.newSnapshotName(prefix))
		err = os.Mkdir(path, os.ModeDir)
	}
	if err != nil {
		path = ""
	}
	return
}

type SnapshotReport struct {
	Keepalive bool        // was the instance alive while running the snapshot?
	Panic     interface{} // was there a panic?

	// errors for the various components of the Snapshot
	StopErr        error
	StartErr       error
	BookkeepingErr error
	FilesystemErr  error
	TriplestoreErr error
	SQLErr         error
}

// Snapshot creates a new snapshot of this instance into dest
func (instance Instance) Snapshot(io stream.IOStream, keepalive bool, dest string) (report SnapshotReport) {
	// catch anything critical that happened during the snapshot
	defer func() {
		report.Panic = recover()
	}()

	stack := instance.Stack()

	// stop the instance (unless it was explicitly asked to not do so!)
	report.Keepalive = keepalive
	if !keepalive {
		logging.LogMessage(io, "Stopping instance")
		report.StopErr = stack.Down(io)
		defer func() {
			logging.LogMessage(io, "Starting instance")
			report.StartErr = stack.Up(io)
		}()
	}

	// create a wait group, and message channel
	wg := &sync.WaitGroup{}
	messages := make(chan string, 4)

	// write bookkeeping information
	wg.Add(1)
	go func() {
		defer wg.Done()

		bkPath := filepath.Join(dest, "bookkeeping.txt")
		messages <- bkPath

		info, err := os.Create(bkPath)
		if err != nil {
			report.BookkeepingErr = err
			return
		}
		defer info.Close()

		// print whatever is in the database
		// TODO: This should be sql code, maybe gorm can do that?
		_, report.BookkeepingErr = fmt.Fprintf(info, "%#v\n", instance.Instance)
	}()

	// backup the filesystem
	wg.Add(1)
	go func() {
		defer wg.Done()

		fsPath := filepath.Join(dest, filepath.Base(instance.FilesystemBase))
		if err := os.Mkdir(fsPath, fs.ModeDir); err != nil {
			report.FilesystemErr = err
			return
		}

		// copy over whatever is in the base directory
		report.FilesystemErr = fsx.CopyDirectory(fsPath, instance.FilesystemBase, func(dst, src string) {
			messages <- dst
		})
	}()

	// backup the graph db repository
	wg.Add(1)
	go func() {
		defer wg.Done()

		tsPath := filepath.Join(dest, instance.GraphDBRepository+".nq")
		messages <- tsPath

		nquads, err := os.Create(tsPath)
		if err != nil {
			report.TriplestoreErr = err
		}
		defer nquads.Close()

		// directly store the result
		_, report.TriplestoreErr = instance.dis.Triplestore().Backup(nquads, instance.GraphDBRepository)
	}()

	// backup the sql database
	wg.Add(1)
	go func() {
		defer wg.Done()

		sqlPath := filepath.Join(dest, instance.SqlDatabase+".sql")
		messages <- sqlPath

		sql, err := os.Create(sqlPath)
		if err != nil {
			report.SQLErr = err
			return
		}
		defer sql.Close()

		// directly store the result
		report.SQLErr = instance.dis.SQL().Backup(io, sql, instance.SqlDatabase)
	}()

	// TODO: Backup the docker image

	// wait for the group, then close the message channel.
	go func() {
		wg.Wait()
		close(messages)
	}()

	// print out all the messages
	for message := range messages {
		io.Println(message)
	}

	return
}
