//spellchecker:words admin
package admin

//spellchecker:words context embed html template http github wisski distillery internal component server assets templating models pkglib httpx julienschmidt httprouter
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
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"go.tkw01536.de/pkglib/httpx"

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

func (admin *Admin) instanceSnapshots(context.Context) http.Handler {
	tpl := instanceSnapshotsTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuInstance,
			menuSnapshots,
		),
	)

	return tpl.HTMLHandlerWithFlags(admin.dependencies.Handling, func(r *http.Request) (ctx instanceSnapshotsContext, funcs []templating.FlagFunc, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		// setup the context with just the instance
		ctx.Instance, err = admin.dependencies.Instances.WissKI(r.Context(), slug)
		if err != nil {
			return ctx, nil, httpx.ErrNotFound
		}

		ctx.Snapshots, err = ctx.Instance.Snapshots(r.Context())
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to get snapshots: %w", err)
		}

		escapedSlug := url.PathEscape(ctx.Instance.Slug)
		return ctx, []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + escapedSlug)}),                  // #nosec G203 -- escaped and safe
			templating.ReplaceCrumb(menuSnapshots, component.MenuItem{Title: "Snapshots", Path: template.URL("/admin/instance/" + escapedSlug + "/snapshots")}), // #nosec G203 -- escaped and safe
			templating.Title(ctx.Instance.Slug + " - Snapshots"),
			admin.instanceTabs(escapedSlug, "snapshots"),
		}, nil
	})
}
