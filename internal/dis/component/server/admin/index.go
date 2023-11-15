package admin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"golang.org/x/sync/errgroup"
)

// Status produces a new observation of the distillery, and a new information of all instances
// The information on all instances is passed the given quick flag.
func (admin *Admin) Status(ctx context.Context, QuickInformation bool) (target status.Distillery, information []status.WissKI, err error) {
	var group errgroup.Group

	group.Go(func() error {
		// list all the instances
		all, err := admin.dependencies.Instances.All(ctx)
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
					if err != nil {
						return fmt.Errorf("instance %q: %w", instance.Slug, err)
					}
					return
				})
			}
		}
		return nil
	})

	// gather all the observations
	flags := component.FetcherFlags{
		Context: ctx,
	}
	for _, o := range admin.dependencies.Fetchers {
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

func (admin *Admin) Fetch(flags component.FetcherFlags, target *status.Distillery) error {
	target.Time = time.Now().UTC()
	target.Config = admin.Config
	return nil
}

//go:embed "html/index.html"
var indexHTML []byte
var indexTemplate = templating.Parse[indexContext](
	"index.html", indexHTML, nil,

	templating.Title("Admin"),
	templating.Assets(assets.AssetsAdmin),

	templating.Crumbs(
		menuAdmin,
	),
)

//go:embed "html/instances.html"
var instancesHTML []byte
var instancesTemplate = templating.Parse[indexContext](
	"instances.html", instancesHTML, nil,

	templating.Title("Instances"),
	templating.Assets(assets.AssetsAdmin),
)

type indexContext struct {
	templating.RuntimeFlags

	status.Distillery
	Instances []status.WissKI
}

func (admin *Admin) index(ctx context.Context) http.Handler {
	tpl := indexTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Actions(
			menuUsers,
			menuInstances,
		),
	)

	return tpl.HTMLHandler(func(r *http.Request) (idx indexContext, err error) {
		idx.Distillery, idx.Instances, err = admin.Status(r.Context(), false)
		return
	})
}

func (admin *Admin) instances(ctx context.Context) http.Handler {
	tpl := instancesTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuInstances,
		),
		templating.Actions(
			menuProvision,
		),
	)

	return tpl.HTMLHandler(func(r *http.Request) (idx indexContext, err error) {
		idx.Distillery, idx.Instances, err = admin.Status(r.Context(), true)
		return
	})
}
