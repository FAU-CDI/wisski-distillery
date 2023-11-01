package admin

import (
	"context"
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/julienschmidt/httprouter"
	"github.com/tkw1536/pkglib/httpx"

	_ "embed"
)

//go:embed "html/instance.html"
var instanceHTML []byte
var instanceTemplate = templating.Parse[instanceContext](
	"instance.html", instanceHTML, nil,

	templating.Assets(assets.AssetsAdmin),
)

type instanceContext struct {
	templating.RuntimeFlags

	Instance models.Instance
	Info     status.WissKI
}

func (admin *Admin) instance(ctx context.Context) http.Handler {
	tpl := instanceTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuInstance,
		),
		templating.Actions(
			menuRebuild,
			menuGrants,
			menuIngredients,
		),
	)

	return tpl.HTMLHandlerWithFlags(func(r *http.Request) (ic instanceContext, funcs []templating.FlagFunc, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		// find the instance itself!
		instance, err := admin.dependencies.Instances.WissKI(r.Context(), slug)
		if err == instances.ErrWissKINotFound {
			return ic, nil, httpx.ErrNotFound
		}
		if err != nil {
			return ic, nil, err
		}
		ic.Instance = instance.Instance

		// get some more info about the wisski
		ic.Info, err = instance.Info().Information(r.Context(), false)
		if err != nil {
			return ic, nil, err
		}

		funcs = []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + slug)}),
			templating.ReplaceAction(menuRebuild, component.MenuItem{Title: "Rebuild", Path: template.URL("/admin/rebuild/" + slug)}),
			templating.ReplaceAction(menuGrants, component.MenuItem{Title: "Grants", Path: template.URL("/admin/grants/" + slug)}),

			templating.Title(instance.Slug),
		}

		return
	})
}
