//spellchecker:words admin
package admin

//spellchecker:words context http time embed github wisski distillery internal component server assets templating status golang sync errgroup
import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"golang.org/x/sync/errgroup"
)

// Status produces a new observation of the distillery, and a new information of all instances
// The information on all instances is passed the given quick flag.
func (admin *Admin) Status(ctx context.Context, quick bool) (target status.Distillery, information []status.WissKI, err error) {
	var group errgroup.Group

	group.Go(func() error {
		// list all the instances
		all, err := admin.dependencies.Instances.All(ctx)
		if err != nil {
			return fmt.Errorf("failed to list all instances: %w", err)
		}

		// get all of their info!
		information = make([]status.WissKI, len(all))

		var wg sync.WaitGroup
		wg.Add(len(all))
		for i, instance := range all {
			go func() {
				defer wg.Done()

				var err error
				information[i], err = instance.Info().Information(ctx, true)

				if err != nil {
					wdlog.Of(ctx).Warn(
						"failed to fetch information for instance",
						"error", err,
						"slug", all[i].Slug,
					)
				}
			}()
		}

		wg.Wait()

		return nil
	})

	// gather all the observations
	flags := component.FetcherFlags{
		Context: ctx,
	}
	for _, o := range admin.dependencies.Fetchers {
		group.Go(func() error {
			return o.Fetch(flags, &target)
		})
	}

	// wait for all the fetchers to finish
	if err := group.Wait(); err != nil {
		return status.Distillery{}, nil, fmt.Errorf("failed to fetch distillery information: %w", err)
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
	target.Config = component.GetStill(admin).Config
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

func (admin *Admin) index(context.Context) http.Handler {
	tpl := indexTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Actions(
			menuUsers,
			menuInstances,
		),
	)

	return tpl.HTMLHandler(admin.dependencies.Handling, func(r *http.Request) (idx indexContext, err error) {
		idx.Distillery, idx.Instances, err = admin.Status(r.Context(), false)
		return
	})
}

func (admin *Admin) instances(context.Context) http.Handler {
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

	return tpl.HTMLHandler(admin.dependencies.Handling, func(r *http.Request) (idx indexContext, err error) {
		idx.Distillery, idx.Instances, err = admin.Status(r.Context(), true)
		return
	})
}
