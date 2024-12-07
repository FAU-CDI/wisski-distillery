//spellchecker:words exporter
package exporter

//spellchecker:words github wisski distillery internal component models
import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

// Filesystem implements snapshotting an instnace filesystem
type Filesystem struct {
	component.Base
}

var (
	_ component.Snapshotable = (*Filesystem)(nil)
)

// SnapshotNeedsRunning returns if this Snapshotable requires a running instance.
func (Filesystem) SnapshotNeedsRunning() bool { return false }

// SnapshotName returns a new name to be used as an argument for path.
func (Filesystem) SnapshotName() string { return "data" }

// Snapshot creates a snapshot of this instance
func (*Filesystem) Snapshot(wisski models.Instance, context *component.StagingContext) error {
	return context.CopyDirectory(".", wisski.FilesystemBase)
}
