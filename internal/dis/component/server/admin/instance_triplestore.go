//spellchecker:words admin
package admin

//spellchecker:words context embed html template http github wisski distillery internal component server assets templating ingredient extras pkglib httpx julienschmidt httprouter
import (
	"context"
	_ "embed"
	"html/template"
	"net/http"
	"net/url"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php/extras"
	"go.tkw01536.de/pkglib/httpx"

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
	Adapters []extras.DistilleryAdapter
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
		ctx.Adapters = ctx.Instance.Adapters().Adapters()

		escapedSlug := url.PathEscape(ctx.Instance.Slug)
		presentFunc, presentErr := admin.preparePanelInstancePage(r, ctx.Instance, "triplestore")
		if presentErr != nil {
			return ctx, nil, presentErr
		}
		return ctx, []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + escapedSlug)}),                        // #nosec G203 -- escaped and safe
			templating.ReplaceCrumb(menuTriplestore, component.MenuItem{Title: "Triplestore", Path: template.URL("/admin/instance/" + escapedSlug + "/triplestore")}), // #nosec G203 -- escaped and safe
			templating.Title(ctx.Instance.Slug + " - Triplestore"),
			presentFunc,
		}, nil
	})
}
