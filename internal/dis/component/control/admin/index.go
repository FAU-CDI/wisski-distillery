package admin

import (
	"context"
	"net/http"
	"time"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static/custom"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"golang.org/x/sync/errgroup"
)

//go:embed "html/index.html"
var indexTemplateStr string
var indexTemplate = static.AssetsAdmin.MustParseShared(
	"index.html",
	indexTemplateStr,
)

// Status produces a new observation of the distillery, and a new information of all instances
// The information on all instances is passed the given quick flag.
func (admin *Admin) Status(ctx context.Context, QuickInformation bool) (target status.Distillery, information []status.WissKI, err error) {
	var group errgroup.Group

	group.Go(func() error {
		// list all the instances
		all, err := admin.Dependencies.Instances.All(ctx)
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
	for _, o := range admin.Dependencies.Fetchers {
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

type indexContext struct {
	custom.BaseContext

	status.Distillery
	Instances []status.WissKI
}

func (admin *Admin) index(r *http.Request) (idx indexContext, err error) {
	admin.Dependencies.Custom.Update(&idx)
	idx.Distillery, idx.Instances, err = admin.Status(r.Context(), true)
	return
}

func (admin *Admin) Fetch(flags component.FetcherFlags, target *status.Distillery) error {
	target.Time = time.Now().UTC()
	target.Config = admin.Config
	return nil
}
