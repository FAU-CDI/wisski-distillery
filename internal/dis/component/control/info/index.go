package info

import (
	"context"
	"net/http"
	"time"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"golang.org/x/sync/errgroup"
)

//go:embed "html/index.html"
var indexTemplateStr string
var indexTemplate = static.AssetsControlIndex.MustParseShared(
	"index.html",
	indexTemplateStr,
)

type indexContext struct {
	status.Distillery
	Instances []status.WissKI
}

func (info *Info) index(r *http.Request) (idx indexContext, err error) {
	idx.Distillery, idx.Instances, err = info.Status(r.Context(), true)
	return
}

// Status produces a new observation of the distillery, and a new information of all instances
// The information on all instances is passed the given quick flag.
func (info *Info) Status(ctx context.Context, QuickInformation bool) (target status.Distillery, information []status.WissKI, err error) {
	var group errgroup.Group

	group.Go(func() error {
		// list all the instances
		all, err := info.Dependencies.Instances.All(ctx)
		if err != nil {
			return err
		}

		// get all of their info!
		information = make([]status.WissKI, len(all))
		for i, instance := range all {
			{
				i := i
				instance := instance

				// store the info for this group!
				group.Go(func() (err error) {
					information[i], err = instance.Info().Information(ctx, true)
					return err
				})
			}
		}
		return nil
	})

	// gather all the observations
	flags := component.FetcherFlags{
		Context: ctx,
	}
	for _, o := range info.Dependencies.Fetchers {
		o := o
		group.Go(func() error {
			return o.Fetch(flags, &target)
		})
	}

	// wait for all the fetchers to finish
	if err := group.Wait(); err != nil {
		return status.Distillery{}, nil, err
	}

	// count overall instances
	for _, i := range information {
		if i.Running {
			target.RunningCount++
		} else {
			target.StoppedCount++
		}
	}
	target.TotalCount = len(information)

	return
}

func (nfo *Info) Fetch(flags component.FetcherFlags, target *status.Distillery) error {
	target.Time = time.Now().UTC()
	target.Config = nfo.Config
	return nil
}
