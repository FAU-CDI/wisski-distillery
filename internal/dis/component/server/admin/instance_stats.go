//spellchecker:words admin
package admin

//spellchecker:words context embed html template http github wisski distillery internal component server assets templating status pkglib httpx julienschmidt httprouter
import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"go.tkw01536.de/pkglib/httpx"

	"github.com/julienschmidt/httprouter"
)

//go:embed "html/instance_stats.html"
var instanceStatsHTML []byte
var instanceStatsTemplate = templating.Parse[instanceStatsContext](
	"instance_stats.html", instanceStatsHTML, nil,

	templating.Assets(assets.AssetsAdmin),
)

type instanceStatsContext struct {
	templating.RuntimeFlags

	Instance   *wisski.WissKI
	Statistics status.Statistics
}

func (admin *Admin) instanceStats(context.Context) http.Handler {
	tpl := instanceStatsTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuInstance,
			menuStats,
		),
	)

	return tpl.HTMLHandlerWithFlags(admin.dependencies.Handling, func(r *http.Request) (ctx instanceStatsContext, funcs []templating.FlagFunc, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		// setup the context with just the instance
		ctx.Instance, err = admin.dependencies.Instances.WissKI(r.Context(), slug)
		if err != nil {
			return ctx, nil, httpx.ErrNotFound
		}

		// read statistics
		ctx.Statistics, err = ctx.Instance.Stats().Get(r.Context(), nil)
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to get stats: %w", err)
		}

		escapedSlug := url.PathEscape(ctx.Instance.Slug)
		presentFunc, presentErr := admin.preparePanelInstancePage(r, ctx.Instance, "stats")
		if presentErr != nil {
			return ctx, nil, presentErr
		}
		return ctx, []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + escapedSlug)}),    // #nosec G203 -- escaped and safe
			templating.ReplaceCrumb(menuStats, component.MenuItem{Title: "SSH", Path: template.URL("/admin/instance/" + escapedSlug + "/stats")}), // #nosec G203 -- escaped and safe
			templating.Title(ctx.Instance.Slug + " - Stats"),
			presentFunc,
		}, nil
	})
}
