//spellchecker:words exporter
package exporter

//spellchecker:words context path filepath time maps slices github wisski distillery internal component models wdlog ingredient locker logging pkglib collection contextx errorsx status
import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"maps"
	"slices"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"go.tkw01536.de/pkglib/collection"
	"go.tkw01536.de/pkglib/contextx"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/status"
)

// SnapshotDescription is a description for a snapshot.
type SnapshotDescription struct {
	Dest      string // destination path
	Keepalive bool   // should we keep the instance alive while making the snapshot?

	Parts []string // SnapshotName()s of the components to include.
}

// Snapshot represents the result of generating a snapshot.
//
//nolint:recvcheck
type Snapshot struct {
	Description SnapshotDescription
	Instance    models.Instance

	// Start and End Time of the snapshot
	StartTime time.Time
	EndTime   time.Time

	// Generic Panic that may have occured
	ErrPanic interface{}
	ErrStart error
	ErrStop  error

	// Errors holds errors for each component
	Errors map[string]error

	// Logs contains logfiles for each component
	Logs map[string]string

	// List of files included
	WithManifest

	// snapshotables that are running and not running
	partsRunning []component.Snapshotable `json:"-"`
	partsStopped []component.Snapshotable `json:"-"`
}

// Snapshot creates a new snapshot of this instance into dest.
func (exporter *Exporter) NewSnapshot(ctx context.Context, instance *wisski.WissKI, progress io.Writer, desc SnapshotDescription) (snapshot Snapshot) {
	// #nosec G104
	logging.LogMessage(progress, "Locking instance") //nolint:errcheck // no way to report error
	if err := instance.Locker().TryLock(ctx); err != nil {
		_, _ = fmt.Fprintln(progress, err)
		_, _ = fmt.Fprintln(progress, "Aborting snapshot creation")

		return Snapshot{
			ErrPanic: err,
		}
	}
	defer func() {
		// #nosec G104
		logging.LogMessage(progress, "Unlocking instance") //nolint:errcheck // no way to report error

		ctx, cancel := contextx.Anyways(ctx, time.Second)
		defer cancel()

		instance.Locker().Unlock(ctx)
	}()

	// setup the snapshot
	snapshot.Description = desc
	exporter.resolveParts(ctx, desc.Parts, &snapshot)
	snapshot.Instance = instance.Instance

	// capture anything critical, and write the end time
	defer func() {
		snapshot.ErrPanic = recover()
	}()

	// do the create keeping track of time!
	_ = logging.LogOperation(func() error {
		snapshot.StartTime = time.Now().UTC()

		wboxerr, wboxmsg := snapshot.makeParts(ctx, progress, exporter, instance, false)
		bboxerr, bboxlog := snapshot.makeParts(ctx, progress, exporter, instance, true)

		snapshot.EndTime = time.Now().UTC()

		// collection all the errors and logs
		snapshot.Errors = collection.Append(wboxerr, bboxerr)
		snapshot.Logs = collection.Append(wboxmsg, bboxlog)

		return nil
	}, progress, "Writing snapshot files")

	slices.Sort(snapshot.Manifest)
	return
}

// resolveParts resolves parts, and writes it into snapshot.Description.Parts.
// Also sets up snapshot.partsRunning and snapshot.partsStopped.
// sends a warning about unknown parts into the logger in context.
func (snapshots *Exporter) resolveParts(ctx context.Context, parts []string, snapshot *Snapshot) {
	partMap := make(map[string]component.Snapshotable, len(snapshots.dependencies.Snapshotable))
	for _, part := range snapshots.dependencies.Snapshotable {
		partMap[part.SnapshotName()] = part
	}

	// filter the parts (if requested)
	if len(parts) != 0 {
		keys := make(map[string]struct{}, len(parts))
		for _, part := range parts {
			keys[part] = struct{}{}
		}

		// delete all the parts which weren't explicitly requested
		for part := range partMap {
			if _, ok := keys[part]; !ok {
				delete(partMap, part)
			} else {
				delete(keys, part)
			}
		}

		// throw a warning for unknown parts
		for key := range keys {
			wdlog.Of(ctx).Warn(
				"ignoring unknown snapshot part",
				"part", key,
			)
		}
	}

	// sort the names of all requested parts
	snapshot.Description.Parts = slices.AppendSeq(make([]string, 0, len(partMap)), maps.Keys(partMap))
	slices.Sort(snapshot.Description.Parts)

	// and setup the map for running and stopped parts!
	for _, name := range snapshot.Description.Parts {
		part := partMap[name]
		if part.SnapshotNeedsRunning() {
			snapshot.partsRunning = append(snapshot.partsRunning, part)
		} else {
			snapshot.partsStopped = append(snapshot.partsStopped, part)
		}
	}
}

func (snapshot *Snapshot) makeParts(ctx context.Context, progress io.Writer, _ *Exporter, instance *wisski.WissKI, needsRunning bool) (errmap map[string]error, logmap map[string]string) {
	if !needsRunning && !snapshot.Description.Keepalive {
		stack, err := instance.Barrel().OpenStack()
		if err == nil {
			defer errorsx.Close(stack, &snapshot.ErrStop, "stack")

			_, _ = logging.LogMessage(progress, "Stopping instance")
			snapshot.ErrStop = stack.Down(ctx, progress)

			defer func() {
				_, _ = logging.LogMessage(progress, "Starting instance")
				snapshot.ErrStart = stack.Start(ctx, progress)
			}()
		} else {
			snapshot.ErrStart = fmt.Errorf("failed to stop stack: %w", err)
		}
	}

	// handle writing the manifest!
	manifest, done := snapshot.handleManifest(snapshot.Description.Dest)
	defer done()

	// create a new status
	st := status.NewWithCompat(progress, 0)
	st.Start()
	defer st.Stop()

	// get the components
	var comps []component.Snapshotable
	if needsRunning {
		comps = snapshot.partsRunning
	} else {
		comps = snapshot.partsStopped
	}

	// run each of the parts
	errors, ids := status.Group[component.Snapshotable, error]{
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

	// keep all the log files
	files := st.Keep()

	// store errors and logs
	errmap = make(map[string]error, len(comps))
	logmap = make(map[string]string, len(comps))

	for i, wc := range comps {
		name := wc.SnapshotName()
		errmap[name] = errors[i]

		// read the logfile
		logfile := files[ids[i]]
		bytes, err := os.ReadFile(logfile) // #nosec G304 -- logfile set dynamically
		if err != nil {
			wdlog.Of(ctx).Error(
				"unable to copy logfile",
				"error", err,
				"component", name,
			)
			continue
		}

		// delete it, but store the content in the results
		if err := os.Remove(logfile); err != nil {
			err = fmt.Errorf("failed to remove logfile: %w", err)
			errmap[name] = errorsx.Combine(errmap[name], err)
		}
		logmap[name] = string(bytes)
	}

	return
}
