package wisski

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/bookkeeping"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/pkg/countwriter"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/FAU-CDI/wisski-distillery/pkg/opgroup"
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
func (snapshot Snapshot) Report(w io.Writer) (int, error) {
	ww := countwriter.NewCountWriter(w)

	// TODO: Errors of the writer!
	encoder := json.NewEncoder(ww)
	encoder.SetIndent("", "  ")

	io.WriteString(ww, "======= Begin Snapshot Report "+snapshot.Instance.Slug+" =======\n")

	fmt.Fprintf(ww, "Slug:  %s\n", snapshot.Instance.Slug)
	fmt.Fprintf(ww, "Dest:  %s\n", snapshot.Description.Dest)

	fmt.Fprintf(ww, "Start: %s\n", snapshot.StartTime)
	fmt.Fprintf(ww, "End:   %s\n", snapshot.EndTime)
	io.WriteString(ww, "\n")

	io.WriteString(ww, "======= Description =======\n")
	encoder.Encode(snapshot.Description)
	io.WriteString(ww, "\n")

	io.WriteString(ww, "======= Instance =======\n")
	encoder.Encode(snapshot.Instance)
	io.WriteString(ww, "\n")

	io.WriteString(ww, "======= Errors =======\n")
	fmt.Fprintf(ww, "Panic:       %v\n", snapshot.ErrPanic)
	fmt.Fprintf(ww, "Start:       %s\n", snapshot.ErrStart)
	fmt.Fprintf(ww, "Stop:        %s\n", snapshot.ErrStop)
	fmt.Fprintf(ww, "Bookkeep:    %s\n", snapshot.ErrBookkeep)
	fmt.Fprintf(ww, "Pathbuilder: %s\n", snapshot.ErrPathbuilder)
	fmt.Fprintf(ww, "Filesystem:  %s\n", snapshot.ErrFilesystem)
	fmt.Fprintf(ww, "Triplestore: %s\n", snapshot.ErrTriplestore)
	fmt.Fprintf(ww, "SQL:         %s\n", snapshot.ErrSQL)
	io.WriteString(ww, "\n")

	io.WriteString(ww, "======= Manifest =======\n")
	for _, file := range snapshot.Manifest {
		io.WriteString(ww, file+"\n")
	}

	io.WriteString(ww, "\n")

	io.WriteString(ww, "======= End Snapshot Report "+snapshot.Instance.Slug+"=======\n")

	return ww.Sum()
}

// Snapshot creates a new snapshot of this instance into dest
func (dis *Distillery) Snapshot(instance instances.WissKI, io stream.IOStream, desc SnapshotDescription) (snapshot Snapshot) {
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

		snapshot.makeBlackbox(io, dis, instance)
		snapshot.makeWhitebox(io, dis, instance)

		snapshot.EndTime = time.Now()
		return nil
	}, io, "Writing snapshot files")

	slices.Sort(snapshot.Manifest)
	return
}

// makeBlackbox runs the blackbox backup of the system.
// It pauses the Instance, if a consistent state is required.
func (snapshot *Snapshot) makeBlackbox(io stream.IOStream, dis *Distillery, instance instances.WissKI) {
	stack := instance.Stack()

	og := opgroup.NewOpGroup[string](4)

	// stop the instance (unless it was explicitly asked to not do so!)
	if !snapshot.Description.Keepalive {
		logging.LogMessage(io, "Stopping instance")
		snapshot.ErrStop = stack.Down(io)

		defer func() {
			logging.LogMessage(io, "Starting instance")
			snapshot.ErrStart = stack.Up(io)
		}()
	}

	// write bookkeeping information
	og.GoErr(func(files chan<- string) error {
		bkPath := filepath.Join(snapshot.Description.Dest, "bookkeeping.txt")
		files <- bkPath

		info, err := os.Create(bkPath)
		if err != nil {
			return err
		}
		defer info.Close()

		// print whatever is in the database
		// TODO: This should be sql code, maybe gorm can do that?
		_, err = fmt.Fprintf(info, "%#v\n", instance.Instance)
		return err
	}, &snapshot.ErrBookkeep)

	// backup the filesystem
	og.GoErr(func(files chan<- string) error {
		fsPath := filepath.Join(snapshot.Description.Dest, filepath.Base(instance.FilesystemBase))

		// copy over whatever is in the base directory
		return fsx.CopyDirectory(fsPath, instance.FilesystemBase, func(dst, src string) {
			files <- dst
		})
	}, &snapshot.ErrFilesystem)

	// backup the graph db repository
	og.GoErr(func(files chan<- string) error {
		tsPath := filepath.Join(snapshot.Description.Dest, instance.GraphDBRepository+".nq")
		files <- tsPath

		nquads, err := os.Create(tsPath)
		if err != nil {
			return err
		}
		defer nquads.Close()

		// directly store the result
		_, err = dis.Triplestore().Backup(nquads, instance.GraphDBRepository)
		return err
	}, &snapshot.ErrTriplestore)

	// backup the sql database
	og.GoErr(func(files chan<- string) error {
		sqlPath := filepath.Join(snapshot.Description.Dest, snapshot.Instance.SqlDatabase+".sql")
		files <- sqlPath

		sql, err := os.Create(sqlPath)
		if err != nil {
			return err
		}
		defer sql.Close()

		// directly store the result
		return dis.SQL().Backup(io, sql, instance.SqlDatabase)
	}, &snapshot.ErrSQL)

	// wait for the group!
	snapshot.waitGroup(io, og)
}

// makeWhitebox runs the whitebox backup of the system.
// The instance should be running during this step.
func (snapshot *Snapshot) makeWhitebox(io stream.IOStream, dis *Distillery, instance instances.WissKI) {
	og := opgroup.NewOpGroup[string](1)

	// write pathbuilders
	og.GoErr(func(files chan<- string) error {

		pbPath := filepath.Join(snapshot.Description.Dest, "pathbuilders")
		files <- pbPath

		// create the directory!
		if err := os.Mkdir(pbPath, fs.ModeDir); err != nil {
			return err
		}

		// put in all the pathbuilders
		return instance.ExportPathbuilders(pbPath)
	}, &snapshot.ErrPathbuilder)

	// wait for the group!
	snapshot.waitGroup(io, og)
}

// waitGroup waits for the
func (snapshot *Snapshot) waitGroup(io stream.IOStream, og *opgroup.OpGroup[string]) {
	// wait for the messages to return
	for file := range og.Wait() {
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
