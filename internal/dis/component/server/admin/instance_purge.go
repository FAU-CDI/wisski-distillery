//spellchecker:words admin
package admin

//spellchecker:words context embed html template http github wisski distillery internal component server assets templating pkglib httpx julienschmidt httprouter
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
	"go.tkw01536.de/pkglib/httpx"

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

		escapedSlug := url.PathEscape(ctx.Instance.Slug)
		presentFunc, presentErr := admin.preparePanelInstancePage(r, ctx.Instance, "purge")
		if presentErr != nil {
			return ctx, nil, presentErr
		}
		return ctx, []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + escapedSlug)}),      // #nosec G203 -- escaped and safe
			templating.ReplaceCrumb(menuPurge, component.MenuItem{Title: "Purge", Path: template.URL("/admin/instance/" + escapedSlug + "/purge")}), // #nosec G203 -- escaped and safe
			templating.Title(ctx.Instance.Slug + " - Purge"),
			presentFunc,
		}, nil
	})
}
