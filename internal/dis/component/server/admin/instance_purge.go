package admin

import (
	"context"
	_ "embed"
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/pkglib/httpx"

	"github.com/julienschmidt/httprouter"
)

//go:embed "html/instance_purge.html"
var instancePurgeHTML []byte
var instancePurgeTemplate = templating.Parse[instancePurgeContext](
	"instance_purge.html", instancePurgeHTML, nil,

	templating.Assets(assets.AssetsAdmin),
)

type instancePurgeContext struct {
	templating.RuntimeFlags

	Instance *wisski.WissKI
}

func (admin *Admin) instancePurge(context.Context) http.Handler {
	tpl := instancePurgeTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuInstance,
			menuPurge,
		),
	)

	return tpl.HTMLHandlerWithFlags(admin.dependencies.Handling, func(r *http.Request) (ctx instancePurgeContext, funcs []templating.FlagFunc, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		// setup the context with just the instance
		ctx.Instance, err = admin.dependencies.Instances.WissKI(r.Context(), slug)
		if err != nil {
			return ctx, nil, httpx.ErrNotFound
		}

		return ctx, []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + ctx.Instance.Slug)}),
			templating.ReplaceCrumb(menuPurge, component.MenuItem{Title: "Purge", Path: template.URL("/admin/instance/" + ctx.Instance.Slug + "/purge")}),
			templating.Title(ctx.Instance.Slug + " - Purge"),
			admin.instanceTabs(slug, "purge"),
		}, nil
	})
}
