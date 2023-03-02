package exporter

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/locker"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/status"
	"github.com/tkw1536/pkglib/collection"
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
func (snapshots *Exporter) NewSnapshot(ctx context.Context, instance *wisski.WissKI, progress io.Writer, desc SnapshotDescription) (snapshot Snapshot) {

	logging.LogMessage(progress, ctx, "Locking instance")
	if !instance.Locker().TryLock(ctx) {
		err := locker.Locked
		logging.ProgressF(progress, ctx, "%v", err)
		logging.LogMessage(progress, ctx, "Aborting snapshot creation")

		return Snapshot{
			ErrPanic: err,
		}
	}
	defer func() {
		logging.LogMessage(progress, ctx, "Unlocking instance")
		instance.Locker().Unlock(ctx)
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

		snapshot.ErrWhitebox = snapshot.makeParts(ctx, progress, snapshots, instance, false)
		snapshot.ErrBlackbox = snapshot.makeParts(ctx, progress, snapshots, instance, true)

		snapshot.EndTime = time.Now().UTC()
		return nil
	}, progress, ctx, "Writing snapshot files")

	slices.Sort(snapshot.Manifest)
	return
}

func (snapshot *Snapshot) makeParts(ctx context.Context, progress io.Writer, snapshots *Exporter, instance *wisski.WissKI, needsRunning bool) map[string]error {
	if !needsRunning && !snapshot.Description.Keepalive {
		stack := instance.Barrel().Stack()

		logging.LogMessage(progress, ctx, "Stopping instance")
		snapshot.ErrStop = stack.Down(ctx, progress)

		defer func() {
			logging.LogMessage(progress, ctx, "Starting instance")
			snapshot.ErrStart = stack.Up(ctx, progress)
		}()
	}
	// handle writing the manifest!
	manifest, done := snapshot.handleManifest(snapshot.Description.Dest)
	defer done()

	// create a new status
	st := status.NewWithCompat(progress, 0)
	st.Start()
	defer st.Stop()

	// get all the components
	comps := collection.FilterClone(snapshots.Dependencies.Snapshotable, func(sc component.Snapshotable) bool {
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
					ctx,
					writer,
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
