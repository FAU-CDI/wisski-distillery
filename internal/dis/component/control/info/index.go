package info

import (
	"net/http"
	"time"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"golang.org/x/sync/errgroup"
)

//go:embed "html/index.html"
var indexTemplateStr string
var indexTemplate = static.AssetsControlIndex.MustParseShared(
	"index.html",
	indexTemplateStr,
)

type indexContext struct {
	component.Observation
	Instances []ingredient.Information
}

func (info *Info) index(r *http.Request) (idx indexContext, err error) {
	idx.Observation, idx.Instances, err = info.Status(true)
	return
}

// Status produces a new observation of the distillery, and a new information of all instances
// The information on all instances is passed the given quick flag.
func (info *Info) Status(QuickInformation bool) (observation component.Observation, information []ingredient.Information, err error) {
	var group errgroup.Group

	group.Go(func() error {
		// list all the instances
		all, err := info.Instances.All()
		if err != nil {
			return err
		}

		// get all of their info!
		information = make([]ingredient.Information, len(all))
		for i, instance := range all {
			{
				i := i
				instance := instance

				// store the info for this group!
				group.Go(func() (err error) {
					information[i], err = instance.Info().Information(true)
					return err
				})
			}
		}
		return nil
	})

	// gather all the observations
	var flags component.ObservationFlags
	for _, o := range info.Obervers {
		o := o
		group.Go(func() error {
			return o.Observe(flags, &observation)
		})
	}

	// wait for all the observes to finish
	if err := group.Wait(); err != nil {
		return component.Observation{}, nil, err
	}

	// count overall instances
	for _, i := range information {
		if i.Running {
			observation.RunningCount++
		} else {
			observation.StoppedCount++
		}
	}
	observation.TotalCount = len(information)

	return
}

func (nfo *Info) Observe(flags component.ObservationFlags, observation *component.Observation) error {
	observation.Time = time.Now().UTC()
	observation.Config = nfo.Config
	return nil
}
