//spellchecker:words triplestore
package triplestore

//spellchecker:words context errors http github wisski distillery internal component models pkglib errorsx
import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

func (Triplestore) SnapshotNeedsRunning(wisski models.Instance) bool { return false }

func (Triplestore) SnapshotName() string { return "triplestore" }

func (ts *Triplestore) Snapshot(wisski models.Instance, scontext *component.StagingContext) error {
	if err := scontext.AddDirectory(".", func(ctx context.Context) error {
		if err := scontext.AddFile(wisski.GraphDBRepository+".nq", func(ctx context.Context, file io.Writer) error {
			_, err := ts.client().ExportContent(ctx, file, wisski.GraphDBRepository)
			if err != nil {
				return fmt.Errorf("failed to snapshot database: %w", err)
			}
			return nil
		}); err != nil {
			return fmt.Errorf("failed to add nq file: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to add directory: %w", err)
	}
	return nil
}
