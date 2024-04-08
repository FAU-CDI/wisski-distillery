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

//go:embed "html/instance_triplestore.html"
var instanceTriplestoreHTML []byte
var instanceTriplestoreTemplate = templating.Parse[instanceTriplestoreContext](
	"instance_triplestore.html", instanceTriplestoreHTML, nil,

	templating.Assets(assets.AssetsAdmin),
)

type instanceTriplestoreContext struct {
	templating.RuntimeFlags

	Instance *wisski.WissKI
}

func (admin *Admin) instanceTS(context.Context) http.Handler {
	tpl := instanceTriplestoreTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuInstance,
			menuTriplestore,
		),
	)

	return tpl.HTMLHandlerWithFlags(admin.dependencies.Handling, func(r *http.Request) (ctx instanceTriplestoreContext, funcs []templating.FlagFunc, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		// setup the context with just the instance
		ctx.Instance, err = admin.dependencies.Instances.WissKI(r.Context(), slug)
		if err != nil {
			return ctx, nil, httpx.ErrNotFound
		}

		return ctx, []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + ctx.Instance.Slug)}),
			templating.ReplaceCrumb(menuTriplestore, component.MenuItem{Title: "Triplestore", Path: template.URL("/admin/instance/" + ctx.Instance.Slug + "/triplestore")}),
			templating.Title(ctx.Instance.Slug + " - Triplestore"),
			admin.instanceTabs(slug, "triplestore"),
		}, nil
	})
}
