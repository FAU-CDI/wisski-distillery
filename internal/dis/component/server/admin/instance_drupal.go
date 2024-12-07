//spellchecker:words admin
package admin

//spellchecker:words context embed html template http time github wisski distillery internal component server assets templating status pkglib httpx golang sync errgroup julienschmidt httprouter
import (
	"context"
	_ "embed"
	"html/template"
	"net/http"
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
		eg.Go(func() (err error) {
			ctx.Requirements, err = ctx.Instance.Requirements().Get(r.Context(), nil)
			return
		})

		// get the drupal version
		eg.Go(func() (err error) {
			ctx.DrupalVersion, err = ctx.Instance.Version().Get(r.Context(), nil)
			return
		})

		// get the default theme
		eg.Go(func() (err error) {
			ctx.DefaultTheme, err = ctx.Instance.Theme().Get(r.Context(), nil)
			return
		})

		// last time cron was executed
		eg.Go(func() (err error) {
			ctx.LastCron, err = ctx.Instance.Drush().LastCron(r.Context(), nil)
			return
		})

		if err = eg.Wait(); err != nil {
			return ctx, nil, err
		}

		return ctx, []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + ctx.Instance.Slug)}),
			templating.ReplaceCrumb(menuDrupal, component.MenuItem{Title: "Drupal", Path: template.URL("/admin/instance/" + ctx.Instance.Slug + "/drupal")}),
			templating.Title(ctx.Instance.Slug + " - Drupal"),
			admin.instanceTabs(slug, "drupal"),
		}, nil
	})
}
