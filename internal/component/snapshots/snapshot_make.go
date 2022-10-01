package snapshots

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/countwriter"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/FAU-CDI/wisski-distillery/pkg/opgroup"
	"github.com/tkw1536/goprogram/status"
	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/exp/slices"
)

// SnapshotDescription is a description for a snapshot
type SnapshotDescription struct {
	Dest      string // destination path
	Log       bool   // should we log the creation of this snapshot?
	Keepalive bool   // should we keep the instance alive while making the snapshot?
}

// Snapshot represents the result of generating a snapshot
type Snapshot struct {
	Description SnapshotDescription
	Instance    models.Instance

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
func (snapshots *Manager) NewSnapshot(instance instances.WissKI, io stream.IOStream, desc SnapshotDescription) (snapshot Snapshot) {
	// setup the snapshot
	snapshot.Description = desc
	snapshot.Instance = instance.Instance

	// capture anything critical, and write the end time
	defer func() {
		snapshot.ErrPanic = recover()
	}()

	// do the create keeping track of time!
	logging.LogOperation(func() error {
		snapshot.StartTime = time.Now().UTC()

		snapshot.makeBlackbox(io, snapshots, instance)
		snapshot.makeWhitebox(io, snapshots, instance)

		snapshot.EndTime = time.Now().UTC()
		return nil
	}, io, "Writing snapshot files")

	slices.Sort(snapshot.Manifest)
	return
}

// makeBlackbox runs the blackbox backup of the system.
// It pauses the Instance, if a consistent state is required.
func (snapshot *Snapshot) makeBlackbox(io stream.IOStream, snapshots *Manager, instance instances.WissKI) {
	stack := instance.Barrel()

	og := opgroup.NewOpGroup[string](4)

	st := status.NewWithCompat(io.Stdout, 0)
	st.Start()
	defer st.Stop()

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
		line := st.OpenLine("[snapshot bookkeeping]: ", "")
		defer line.Close()
		defer fmt.Fprintln(line, "done")

		bkPath := filepath.Join(snapshot.Description.Dest, "bookkeeping.txt")
		fmt.Fprintln(line, bkPath)
		files <- bkPath

		info, err := snapshots.Core.Environment.Create(bkPath, environment.DefaultFilePerm)
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
		line := st.OpenLine("[snapshot filesystem]: ", "")
		defer line.Close()
		defer fmt.Fprintln(line, "done")

		fsPath := filepath.Join(snapshot.Description.Dest, filepath.Base(instance.FilesystemBase))

		// copy over whatever is in the base directory
		defer fmt.Fprintln(line, "done")
		return fsx.CopyDirectory(snapshots.Core.Environment, fsPath, instance.FilesystemBase, func(dst, src string) {
			fmt.Fprintln(line, dst)
			files <- dst
		})

	}, &snapshot.ErrFilesystem)

	// backup the graph db repository
	og.GoErr(func(files chan<- string) error {
		line := st.OpenLine("[snapshot triplestore]: ", "")
		defer line.Close()
		defer fmt.Fprintln(line, "done")

		tsPath := filepath.Join(snapshot.Description.Dest, instance.GraphDBRepository+".nq")
		fmt.Fprintln(line, tsPath)
		files <- tsPath

		nquads, err := snapshots.Core.Environment.Create(tsPath, environment.DefaultFilePerm)
		if err != nil {
			return err
		}
		defer nquads.Close()

		// directly store the result
		_, err = snapshots.TS.Snapshot(nquads, instance.GraphDBRepository)
		return err
	}, &snapshot.ErrTriplestore)

	// backup the sql database
	og.GoErr(func(files chan<- string) error {
		line := st.OpenLine("[snapshot sql]: ", "")
		defer line.Close()
		defer fmt.Fprintln(line, "done")

		sqlPath := filepath.Join(snapshot.Description.Dest, snapshot.Instance.SqlDatabase+".sql")
		fmt.Fprintln(line, sqlPath)
		files <- sqlPath

		sql, err := snapshots.Core.Environment.Create(sqlPath, environment.DefaultFilePerm)
		if err != nil {
			return err
		}
		defer sql.Close()

		// directly store the result
		return snapshots.SQL.Snapshot(io, sql, instance.SqlDatabase)
	}, &snapshot.ErrSQL)

	// wait for the group!
	snapshot.waitGroup(io, og)
}

// makeWhitebox runs the whitebox backup of the system.
// The instance should be running during this step.
func (snapshot *Snapshot) makeWhitebox(io stream.IOStream, snapshots *Manager, instance instances.WissKI) {
	og := opgroup.NewOpGroup[string](1)

	// write pathbuilders
	og.GoErr(func(files chan<- string) error {

		pbPath := filepath.Join(snapshot.Description.Dest, "pathbuilders")
		files <- pbPath

		// create the directory!
		if err := snapshots.Core.Environment.Mkdir(pbPath, environment.DefaultDirPerm); err != nil {
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

		// add the manifest
		snapshot.Manifest = append(snapshot.Manifest, path)
	}
}

// WriteReport writes out the report belonging to this snapshot.
// It is a separate function, to allow writing it indepenently of the rest.
func (snapshot *Snapshot) WriteReport(env environment.Environment, stream stream.IOStream) error {
	return logging.LogOperation(func() error {
		reportPath := filepath.Join(snapshot.Description.Dest, "report.txt")
		stream.Println(reportPath)

		// create the report file!
		report, err := env.Create(reportPath, environment.DefaultFilePerm)
		if err != nil {
			return err
		}
		defer report.Close()

		// print the report into it!
		_, err = io.WriteString(report, snapshot.String())
		return err
	}, stream, "Writing snapshot report")
}
