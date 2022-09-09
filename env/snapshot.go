package env

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/bookkeeping"
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
	if prefix == "" {
		prefix = "backup"
	} else {
		prefix = "snapshot-" + prefix
	}
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

// SnapshotDescription is a description for a snapshot
type SnapshotDescription struct {
	Dest      string // destination path
	Keepalive bool   // should we keep the instance alive while making the snapshot?
}

// Snapshot represents the result of generating a snapshot
type Snapshot struct {
	Description SnapshotDescription
	Instance    bookkeeping.Instance

	// various error states, which are ignored when creating the snapshot
	ErrPanic interface{} // panic, if any

	ErrStart error
	ErrStop  error

	ErrBookkeep    error
	ErrPathbuilder error
	ErrFilesystem  error
	ErrTriplestore error
	ErrSSQL        error
}

// Snapshot creates a new snapshot of this instance into dest
func (instance Instance) Snapshot(io stream.IOStream, desc SnapshotDescription) (snapshot Snapshot) {
	// setup the snapshot
	snapshot.Description = desc
	snapshot.Instance = instance.Instance

	// catch anything critical that happened during the snapshot
	defer func() {
		snapshot.ErrPanic = recover()
	}()

	// and do the create!
	snapshot.create(io, instance)

	return
}

func (snapshot *Snapshot) create(io stream.IOStream, instance Instance) {
	stack := instance.Stack()

	// stop the instance (unless it was explicitly asked to not do so!)
	if !snapshot.Description.Keepalive {
		logging.LogMessage(io, "Stopping instance")
		snapshot.ErrStop = stack.Down(io)
		defer func() {
			logging.LogMessage(io, "Starting instance")
			snapshot.ErrStart = stack.Up(io)
		}()
	}

	// create a wait group, and message channel
	wg := &sync.WaitGroup{}
	messages := make(chan string, 4)

	// write bookkeeping information
	wg.Add(1)
	go func() {
		defer wg.Done()

		bkPath := filepath.Join(snapshot.Description.Dest, "bookkeeping.txt")
		messages <- bkPath

		info, err := os.Create(bkPath)
		if err != nil {
			snapshot.ErrBookkeep = err
			return
		}
		defer info.Close()

		// print whatever is in the database
		// TODO: This should be sql code, maybe gorm can do that?
		_, snapshot.ErrBookkeep = fmt.Fprintf(info, "%#v\n", instance.Instance)
	}()

	// write pathbuilders
	wg.Add(1)
	go func() {
		defer wg.Done()

		pbPath := filepath.Join(snapshot.Description.Dest, "pathbuilders")
		messages <- pbPath

		// create the directory!
		if err := os.Mkdir(pbPath, fs.ModeDir); err != nil {
			snapshot.ErrPathbuilder = err
			return
		}

		// put in all the pathbuilders
		snapshot.ErrPathbuilder = instance.ExportPathbuilders(pbPath)
	}()

	// backup the filesystem
	wg.Add(1)
	go func() {
		defer wg.Done()

		fsPath := filepath.Join(snapshot.Description.Dest, filepath.Base(instance.FilesystemBase))
		if err := os.Mkdir(fsPath, fs.ModeDir); err != nil {
			snapshot.ErrFilesystem = err
			return
		}

		// copy over whatever is in the base directory
		snapshot.ErrFilesystem = fsx.CopyDirectory(fsPath, instance.FilesystemBase, func(dst, src string) {
			messages <- dst
		})
	}()

	// backup the graph db repository
	wg.Add(1)
	go func() {
		defer wg.Done()

		tsPath := filepath.Join(snapshot.Description.Dest, instance.GraphDBRepository+".nq")
		messages <- tsPath

		nquads, err := os.Create(tsPath)
		if err != nil {
			snapshot.ErrTriplestore = err
		}
		defer nquads.Close()

		// directly store the result
		_, snapshot.ErrTriplestore = instance.dis.Triplestore().Backup(nquads, instance.GraphDBRepository)
	}()

	// backup the sql database
	wg.Add(1)
	go func() {
		defer wg.Done()

		sqlPath := filepath.Join(snapshot.Description.Dest, snapshot.Instance.SqlDatabase+".sql")
		messages <- sqlPath

		sql, err := os.Create(sqlPath)
		if err != nil {
			snapshot.ErrSSQL = err
			return
		}
		defer sql.Close()

		// directly store the result
		snapshot.ErrSSQL = instance.dis.SQL().Backup(io, sql, instance.SqlDatabase)
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
}

// WriteReport writes out the report belonging to this snapshot.
// It is a separate function, to allow writing it indepenently of the rest.
func (snapshot Snapshot) WriteReport(io stream.IOStream) error {
	return logging.LogOperation(func() error {
		reportPath := filepath.Join(snapshot.Description.Dest, "report.txt")
		io.Println(reportPath)

		// create the report file!
		report, err := os.Create(reportPath)
		if err != nil {
			return err
		}
		defer report.Close()

		// print the report into it!
		_, err = fmt.Fprintf(report, "%#v\n", snapshot)
		return err
	}, io, "Writing snapshot report")
}
