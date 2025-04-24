//spellchecker:words admin
package admin

//spellchecker:words context embed html template http time github wisski distillery internal component server assets templating status pkglib httpx golang sync errgroup julienschmidt httprouter
import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/pkglib/httpx"
	"golang.org/x/sync/errgroup"

	"github.com/julienschmidt/httprouter"
)

//go:embed "html/instance_drupal.html"
var instanceDrupalHTML []byte
var instanceDrupalTemplate = templating.Parse[instanceDrupalContext](
	"instance_drupal.html", instanceDrupalHTML, nil,

	templating.Assets(assets.AssetsAdmin),
)

type instanceDrupalContext struct {
	templating.RuntimeFlags

	Instance *wisski.WissKI

	DrupalVersion string
	DefaultTheme  string

	Requirements []status.Requirement

	LastCron time.Time
}

func (admin *Admin) instanceDrupal(context.Context) http.Handler {
	tpl := instanceDrupalTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuInstance,
			menuDrupal,
		),
	)

	return tpl.HTMLHandlerWithFlags(admin.dependencies.Handling, func(r *http.Request) (ctx instanceDrupalContext, funcs []templating.FlagFunc, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		// setup the context with just the instance
		ctx.Instance, err = admin.dependencies.Instances.WissKI(r.Context(), slug)
		if err != nil {
			return ctx, nil, httpx.ErrNotFound
		}

		var eg errgroup.Group

		// get the requirements
		//nolint:contextcheck
		eg.Go(func() (err error) {
			ctx.Requirements, err = ctx.Instance.Requirements().Get(r.Context(), nil)
			return
		})

		// get the drupal version
		//nolint:contextcheck
		eg.Go(func() (err error) {
			ctx.DrupalVersion, err = ctx.Instance.Version().Get(r.Context(), nil)
			return
		})

		// get the default theme
		//nolint:contextcheck
		eg.Go(func() (err error) {
			ctx.DefaultTheme, err = ctx.Instance.Theme().Get(r.Context(), nil)
			return
		})

		// last time cron was executed
		//nolint:contextcheck
		eg.Go(func() (err error) {
			ctx.LastCron, err = ctx.Instance.Drush().LastCron(r.Context(), nil)
			return
		})

		if err = eg.Wait(); err != nil {
			return ctx, nil, fmt.Errorf("failed to get values: %w", err)
		}

		escapedSlug := url.PathEscape(ctx.Instance.Slug)
		return ctx, []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + escapedSlug)}),         // #nosec G203 -- escaped and safe
			templating.ReplaceCrumb(menuDrupal, component.MenuItem{Title: "Drupal", Path: template.URL("/admin/instance/" + escapedSlug + "/drupal")}), // #nosec G203 -- escaped and safe
			templating.Title(ctx.Instance.Slug + " - Drupal"),
			admin.instanceTabs(escapedSlug, "drupal"),
		}, nil
	})
}
