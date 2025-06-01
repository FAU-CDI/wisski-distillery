//spellchecker:words admin
package admin

//spellchecker:words context embed html template http github wisski distillery internal component server assets templating pkglib httpx julienschmidt httprouter
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
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php/extras"
	"github.com/tkw1536/pkglib/httpx"

	"github.com/julienschmidt/httprouter"
)

//go:embed "html/instance_modules.html"
var instanceModulesHTML []byte
var instanceModulesTemplate = templating.Parse[instanceModulesContext](
	"instance_purge.html", instanceModulesHTML, nil,

	templating.Assets(assets.AssetsAdmin),
)

type instanceModulesContext struct {
	templating.RuntimeFlags

	Instance *wisski.WissKI
	Modules  []extras.DrushExtendedModuleInfo

	EnabledCount        int
	DisabledCount       int
	CustomEnabledCount  int
	CustomDisabledCount int
}

func (admin *Admin) instanceModules(context.Context) http.Handler {
	tpl := instanceModulesTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuInstance,
			menuPurge,
		),
	)

	return tpl.HTMLHandlerWithFlags(admin.dependencies.Handling, func(r *http.Request) (ctx instanceModulesContext, funcs []templating.FlagFunc, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		// setup the context with just the instance
		ctx.Instance, err = admin.dependencies.Instances.WissKI(r.Context(), slug)
		if err != nil {
			return ctx, nil, httpx.ErrNotFound
		}

		// get all the modules
		ctx.Modules, err = ctx.Instance.Modules().Get(r.Context(), nil)
		if err != nil {
			return ctx, nil, fmt.Errorf("%w: failed to get modules: %w", httpx.ErrInternalServerError, err)
		}

		for _, m := range ctx.Modules {
			if m.Enabled {
				ctx.EnabledCount++
			} else {
				ctx.DisabledCount++
			}
			if m.HasComposer() {
				continue
			}

			if m.Enabled {
				ctx.CustomEnabledCount++
			} else {
				ctx.CustomDisabledCount++
			}
		}

		escapedSlug := url.PathEscape(ctx.Instance.Slug)
		return ctx, []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + escapedSlug)}),          // #nosec G203 -- escaped and safe
			templating.ReplaceCrumb(menuPurge, component.MenuItem{Title: "Modules", Path: template.URL("/admin/instance/" + escapedSlug + "/modules")}), // #nosec G203 -- escaped and safe
			templating.Title(ctx.Instance.Slug + " - Modules"),
			admin.instanceTabs(escapedSlug, "modules"),
		}, nil
	})
}
