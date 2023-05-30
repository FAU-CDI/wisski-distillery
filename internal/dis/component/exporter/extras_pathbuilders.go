package exporter

import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

type Pathbuilders struct {
	component.Base
	Dependencies struct {
		Instances *instances.Instances
	}
}

var (
	_ component.Snapshotable = (*Pathbuilders)(nil)
)

func (Pathbuilders) SnapshotNeedsRunning() bool { return true }

func (Pathbuilders) SnapshotName() string { return "pathbuilders" }

func (pbs *Pathbuilders) Snapshot(wisski models.Instance, scontext *component.StagingContext) error {
	return scontext.AddDirectory(".", func(ctx context.Context) error {
		builders, err := pbs.Dependencies.Instances.Instance(ctx, wisski).Pathbuilder().GetAll(ctx, nil)
		if err != nil {
			return err
		}

		for name, bytes := range builders {
			if err := scontext.AddFile(name+".xml", func(ctx context.Context, file io.Writer) error {
				_, err := file.Write([]byte(bytes))
				return err
			}); err != nil {
				return err
			}
		}
		return nil
	})
}
