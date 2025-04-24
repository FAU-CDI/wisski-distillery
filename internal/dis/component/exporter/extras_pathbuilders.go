//spellchecker:words exporter
package exporter

//spellchecker:words context github wisski distillery internal component instances models
import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

//nolint:recvcheck
type Pathbuilders struct {
	component.Base
	dependencies struct {
		Instances *instances.Instances
	}
}

var (
	_ component.Snapshotable = (*Pathbuilders)(nil)
)

func (Pathbuilders) SnapshotNeedsRunning() bool { return true }

func (Pathbuilders) SnapshotName() string { return "pathbuilders" }

func (pbs *Pathbuilders) Snapshot(wisski models.Instance, scontext *component.StagingContext) error {
	if err := scontext.AddDirectory(".", func(ctx context.Context) error {
		builders, err := pbs.dependencies.Instances.Instance(ctx, wisski).Pathbuilder().GetAll(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to get pathbuilders: %w", err)
		}

		for name, bytes := range builders {
			if err := scontext.AddFile(name+".xml", func(ctx context.Context, file io.Writer) error {
				_, err := file.Write([]byte(bytes))
				if err != nil {
					return fmt.Errorf("failed to write file: %w", err)
				}
				return nil
			}); err != nil {
				return fmt.Errorf("failed to export pathbuilder %q: %w", name, err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to add directory: %w", err)
	}
	return nil
}
