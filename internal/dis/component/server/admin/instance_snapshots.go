package admin

import (
	"context"
	_ "embed"
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/pkglib/httpx"

	"github.com/julienschmidt/httprouter"
)

//go:embed "html/instance_snapshots.html"
var instanceSnapshotsHTML []byte
var instanceSnapshotsTemplate = templating.Parse[instanceSnapshotsContext](
	"instance_snapshots.html", instanceSnapshotsHTML, nil,

	templating.Assets(assets.AssetsAdmin),
)

type instanceSnapshotsContext struct {
	templating.RuntimeFlags

	Instance  *wisski.WissKI
	Snapshots []models.Export
}

func (admin *Admin) instanceSnapshots(ctx context.Context) http.Handler {
	tpl := instanceSnapshotsTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuInstance,
			menuSnapshots,
		),
	)

	return tpl.HTMLHandlerWithFlags(func(r *http.Request) (ctx instanceSnapshotsContext, funcs []templating.FlagFunc, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		// setup the context with just the instance
		ctx.Instance, err = admin.dependencies.Instances.WissKI(r.Context(), slug)
		if err != nil {
			return ctx, nil, httpx.ErrNotFound
		}

		ctx.Snapshots, err = ctx.Instance.Snapshots(r.Context())
		if err != nil {
			return ctx, nil, httpx.ErrInternalServerError
		}

		return ctx, []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + ctx.Instance.Slug)}),
			templating.ReplaceCrumb(menuSnapshots, component.MenuItem{Title: "Snapshots", Path: template.URL("/admin/instance/" + ctx.Instance.Slug + "/snapshots")}),
			templating.Title(ctx.Instance.Slug + " - Snapshots"),
			admin.instanceTabs(slug, "snapshots"),
		}, nil
	})
}
