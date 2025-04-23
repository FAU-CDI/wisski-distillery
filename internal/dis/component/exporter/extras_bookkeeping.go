//spellchecker:words exporter
package exporter

//spellchecker:words context github wisski distillery internal component models
import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

//nolint:recvcheck
type Bookkeeping struct {
	component.Base
}

var (
	_ component.Snapshotable = (*Bookkeeping)(nil)
)

// SnapshotNeedsRunning returns if this Snapshotable requires a running instance.
func (Bookkeeping) SnapshotNeedsRunning() bool { return false }

// SnapshotName returns a new name to be used as an argument for path.
func (Bookkeeping) SnapshotName() string { return "bookkeeping.txt" }

// Snapshot creates a snapshot of this instance.
func (*Bookkeeping) Snapshot(wisski models.Instance, scontext *component.StagingContext) error {
	return scontext.AddFile(".", func(ctx context.Context, file io.Writer) error {
		_, err := fmt.Fprintf(file, "%#v\n", wisski)
		return err
	})
}
