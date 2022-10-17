package snapshots

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/lib/collection"
	"github.com/tkw1536/goprogram/status"
	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/exp/slices"
)

// SnapshotDescription is a description for a snapshot
type SnapshotDescription struct {
	Dest      string // destination path
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
	ErrPanic    interface{}
	ErrStart    error
	ErrStop     error
	ErrWhitebox map[string]error
	ErrBlackbox map[string]error

	// List of files included
	WithManifest
}

// Snapshot creates a new snapshot of this instance into dest
func (snapshots *Manager) NewSnapshot(instance *wisski.WissKI, io stream.IOStream, desc SnapshotDescription) (snapshot Snapshot) {

	logging.LogMessage(io, "Locking instance")
	if err := instance.TryLock(); err != nil {
		io.EPrintln(err)
		logging.LogMessage(io, "Aborting snapshot creation")

		return Snapshot{
			ErrPanic: err,
		}
	}
	defer func() {
		logging.LogMessage(io, "Unlocking instance")
		instance.Unlock()
	}()

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

		snapshot.ErrWhitebox = snapshot.makeParts(io, snapshots, instance, false)
		snapshot.ErrBlackbox = snapshot.makeParts(io, snapshots, instance, true)

		snapshot.EndTime = time.Now().UTC()
		return nil
	}, io, "Writing snapshot files")

	slices.Sort(snapshot.Manifest)
	return
}

func (snapshot *Snapshot) makeParts(ios stream.IOStream, snapshots *Manager, instance *wisski.WissKI, needsRunning bool) map[string]error {
	if !needsRunning && !snapshot.Description.Keepalive {
		stack := instance.Barrel()

		logging.LogMessage(ios, "Stopping instance")
		snapshot.ErrStop = stack.Down(ios)

		defer func() {
			logging.LogMessage(ios, "Starting instance")
			snapshot.ErrStart = stack.Up(ios)
		}()
	}
	// handle writing the manifest!
	manifest, done := snapshot.handleManifest(snapshot.Description.Dest)
	defer done()

	// create a new status
	st := status.NewWithCompat(ios.Stdout, 0)
	st.Start()
	defer st.Stop()

	// get all the components
	comps := collection.FilterClone(snapshots.Snapshotable, func(sc component.Snapshotable) bool {
		return sc.SnapshotNeedsRunning() == needsRunning
	})

	results := make(map[string]error, len(comps))

	errors := status.Group[component.Snapshotable, error]{
		PrefixString: func(item component.Snapshotable, index int) string {
			return fmt.Sprintf("[snapshot %q]: ", item.Name())
		},
		PrefixAlign: true,

		Handler: func(sc component.Snapshotable, index int, writer io.Writer) error {
			return sc.Snapshot(
				instance.Instance,
				component.NewStagingContext(
					snapshots.Environment,
					stream.NewIOStream(writer, writer, nil, 0),
					filepath.Join(snapshot.Description.Dest, sc.SnapshotName()),
					manifest,
				),
			)
		},

		ResultString: status.DefaultErrorString[component.Snapshotable],
	}.Use(st, comps)

	for i, wc := range comps {
		results[wc.Name()] = errors[i]
	}
	return results
}
