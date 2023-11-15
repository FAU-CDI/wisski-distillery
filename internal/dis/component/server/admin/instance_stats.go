package admin

import (
	"context"
	_ "embed"
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/pkglib/httpx"

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

func (admin *Admin) instanceStats(ctx context.Context) http.Handler {
	tpl := instanceStatsTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuInstance,
			menuStats,
		),
	)

	return tpl.HTMLHandlerWithFlags(func(r *http.Request) (ctx instanceStatsContext, funcs []templating.FlagFunc, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		// setup the context with just the instance
		ctx.Instance, err = admin.dependencies.Instances.WissKI(r.Context(), slug)
		if err != nil {
			return ctx, nil, httpx.ErrNotFound
		}

		// read statistics
		ctx.Statistics, err = ctx.Instance.Stats().Get(r.Context(), nil)
		if err != nil {
			return ctx, nil, err
		}

		return ctx, []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + ctx.Instance.Slug)}),
			templating.ReplaceCrumb(menuStats, component.MenuItem{Title: "SSH", Path: template.URL("/admin/instance/" + ctx.Instance.Slug + "/stats")}),
			templating.Title(ctx.Instance.Slug + " - Stats"),
			admin.instanceTabs(slug, "stats"),
		}, nil
	})
}
