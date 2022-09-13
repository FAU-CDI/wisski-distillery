package env

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

	"github.com/FAU-CDI/wisski-distillery/pkg/bookkeeping"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/FAU-CDI/wisski-distillery/pkg/password"
	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/exp/slices"
)

// SnapshotsDir returns the path that contains all snapshot related data.
func (dis *Distillery) SnapshotsDir() string {
	return filepath.Join(dis.Config.DeployRoot, "snapshots")
}

// SnapshotsStagingPath returns the path to the directory containing a temporary staging area for snapshots.
// Use NewSnapshotStagingDir to generate a new staging area.
func (dis *Distillery) SnapshotsStagingPath() string {
	return filepath.Join(dis.SnapshotsDir(), "staging")
}

// SnapshotsArchivePath returns the path to the directory containing all exported archives.
// Use NewSnapshotArchivePath to generate a path to a new archive in this directory.
func (dis *Distillery) SnapshotsArchivePath() string {
	return filepath.Join(dis.SnapshotsDir(), "archives")
}

// NewSnapshotArchivePath returns the path to a new archive with the provided prefix.
// The path is guaranteed to not exist.
func (dis *Distillery) NewSnapshotArchivePath(prefix string) (path string) {
	// TODO: Consider moving these into a subdirectory with the provided prefix.
	for path == "" || fsx.Exists(path) {
		name := dis.newSnapshotName(prefix) + ".tar.gz"
		path = filepath.Join(dis.SnapshotsArchivePath(), name)
	}
	return
}

// newSnapshot name returns a new basename for a snapshot with the provided prefix.
// The name is guaranteed to be unique within this process.
func (*Distillery) newSnapshotName(prefix string) string {
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
func (dis *Distillery) NewSnapshotStagingDir(prefix string) (path string, err error) {
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

	// Start and End Time of the snapshot
	StartTime time.Time
	EndTime   time.Time

	// Generic Panic that may have occured
	ErrPanic interface{}

	// Errors during starting and stopping the system
	ErrStart error
	ErrStop  error

	// List of files included
	Manifest []string

	// Errors during other parts
	ErrBookkeep    error
	ErrPathbuilder error
	ErrFilesystem  error
	ErrTriplestore error
	ErrSQL         error
}

func (snapshot Snapshot) String() string {
	var builder strings.Builder
	snapshot.Report(&builder)
	return builder.String()
}

// Report writes a report from snapshot into w
func (snapshot Snapshot) Report(w io.Writer) {
	// TODO: Errors of the writer!
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	io.WriteString(w, "======= Begin Snapshot Report "+snapshot.Instance.Slug+" =======\n")

	fmt.Fprintf(w, "Slug:  %s\n", snapshot.Instance.Slug)
	fmt.Fprintf(w, "Dest:  %s\n", snapshot.Description.Dest)

	fmt.Fprintf(w, "Start: %s\n", snapshot.StartTime)
	fmt.Fprintf(w, "End:   %s\n", snapshot.EndTime)
	io.WriteString(w, "\n")

	io.WriteString(w, "======= Description =======\n")
	encoder.Encode(snapshot.Description)
	io.WriteString(w, "\n")

	io.WriteString(w, "======= Instance =======\n")
	encoder.Encode(snapshot.Instance)
	io.WriteString(w, "\n")

	io.WriteString(w, "======= Errors =======\n")
	fmt.Fprintf(w, "Panic:       %v\n", snapshot.ErrPanic)
	fmt.Fprintf(w, "Start:       %s\n", snapshot.ErrStart)
	fmt.Fprintf(w, "Stop:        %s\n", snapshot.ErrStop)
	fmt.Fprintf(w, "Bookkeep:    %s\n", snapshot.ErrBookkeep)
	fmt.Fprintf(w, "Pathbuilder: %s\n", snapshot.ErrPathbuilder)
	fmt.Fprintf(w, "Filesystem:  %s\n", snapshot.ErrFilesystem)
	fmt.Fprintf(w, "Triplestore: %s\n", snapshot.ErrTriplestore)
	fmt.Fprintf(w, "SQL:         %s\n", snapshot.ErrSQL)
	io.WriteString(w, "\n")

	io.WriteString(w, "======= Manifest =======\n")
	for _, file := range snapshot.Manifest {
		io.WriteString(w, file+"\n")
	}

	io.WriteString(w, "\n")

	io.WriteString(w, "======= End Snapshot Report "+snapshot.Instance.Slug+" =======\n")
}

// Snapshot creates a new snapshot of this instance into dest
func (instance Instance) Snapshot(io stream.IOStream, desc SnapshotDescription) (snapshot Snapshot) {
	// setup the snapshot
	snapshot.Description = desc
	snapshot.Instance = instance.Instance

	// capture anything critical, and write the end time
	defer func() {
		snapshot.ErrPanic = recover()
	}()

	// do the create keeping track of time!
	logging.LogOperation(func() error {
		snapshot.StartTime = time.Now()
		snapshot.create(io, instance)
		snapshot.EndTime = time.Now()

		return nil
	}, io, "Writing snapshot files")

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
	files := make(chan string, 4)

	// write bookkeeping information
	wg.Add(1)
	go func() {
		defer wg.Done()

		bkPath := filepath.Join(snapshot.Description.Dest, "bookkeeping.txt")
		files <- bkPath

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
	// TODO: Move this outside of the up/down stuff!
	wg.Add(1)
	go func() {
		defer wg.Done()

		pbPath := filepath.Join(snapshot.Description.Dest, "pathbuilders")
		files <- pbPath

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

		// copy over whatever is in the base directory
		snapshot.ErrFilesystem = fsx.CopyDirectory(fsPath, instance.FilesystemBase, func(dst, src string) {
			files <- dst
		})
	}()

	// backup the graph db repository
	wg.Add(1)
	go func() {
		defer wg.Done()

		tsPath := filepath.Join(snapshot.Description.Dest, instance.GraphDBRepository+".nq")
		files <- tsPath

		nquads, err := os.Create(tsPath)
		if err != nil {
			snapshot.ErrTriplestore = err
			return
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
		files <- sqlPath

		sql, err := os.Create(sqlPath)
		if err != nil {
			snapshot.ErrSQL = err
			return
		}
		defer sql.Close()

		// directly store the result
		snapshot.ErrSQL = instance.dis.SQL().Backup(io, sql, instance.SqlDatabase)
	}()

	// TODO: Backup the docker image

	// wait for the group, then close the message channel.
	go func() {
		wg.Wait()
		close(files)
	}()

	for file := range files {
		// get the relative path to the root of the manifest.
		// nothing *should* go wrong, but in case it does, use the original path.
		path, err := filepath.Rel(snapshot.Description.Dest, file)
		if err != nil {
			path = file
		}

		// write it to the command line
		// and also add it to the manifest
		io.Printf("\033[2K\r%s", path)
		snapshot.Manifest = append(snapshot.Manifest, path)
	}
	io.Println("")

	// make sure the manifest is sorted!
	slices.Sort(snapshot.Manifest)
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
		_, err = report.WriteString(snapshot.String())
		return err
	}, io, "Writing snapshot report")
}
