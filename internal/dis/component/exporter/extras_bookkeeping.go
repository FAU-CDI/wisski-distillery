package exporter

import (
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

type Bookkeeping struct {
	component.Base
}

// SnapshotNeedsRunning returns if this Snapshotable requires a running instance.
func (Bookkeeping) SnapshotNeedsRunning() bool { return false }

// SnapshotName returns a new name to be used as an argument for path.
func (Bookkeeping) SnapshotName() string { return "bookkeeping.txt" }

// Snapshot creates a snapshot of this instance
func (*Bookkeeping) Snapshot(wisski models.Instance, context component.StagingContext) error {
	return context.AddFile(".", func(file io.Writer) error {
		_, err := fmt.Fprintf(file, "%#v\n", wisski)
		return err
	})
}
