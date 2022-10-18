package exporter

import (
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

type Pathbuilders struct {
	component.Base
	Instances *instances.Instances
}

func (Pathbuilders) SnapshotNeedsRunning() bool { return true }

func (Pathbuilders) SnapshotName() string { return "pathbuilders" }

func (pbs *Pathbuilders) Snapshot(wisski models.Instance, context component.StagingContext) error {
	return context.AddDirectory(".", func() error {
		builders, err := pbs.Instances.Instance(wisski).AllPathbuilders(nil)
		if err != nil {
			return err
		}

		for name, bytes := range builders {
			if err := context.AddFile(name+".xml", func(file io.Writer) error {
				_, err := file.Write([]byte(bytes))
				return err
			}); err != nil {
				return err
			}
		}
		return nil
	})
}
