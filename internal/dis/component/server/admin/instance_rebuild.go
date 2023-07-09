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
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/julienschmidt/httprouter"
	"github.com/tkw1536/pkglib/httpx"

	_ "embed"
)

//go:embed "html/instance_rebuild.html"
var instanceRebuildHTML []byte
var instanceRebuildTemplate = templating.Parse[instanceRebuildContext](
	"instance_rebuild.html", instanceRebuildHTML, nil,

	templating.Title("Rebuild Instance"),
	templating.Assets(assets.AssetsAdminRebuild),
)

type instanceRebuildContext struct {
	templating.RuntimeFlags

	Slug   string
	System models.System

	systemParams
}

func (admin *Admin) instanceRebuild(ctx context.Context) http.Handler {
	tpl := instanceRebuildTemplate.Prepare(
		admin.Dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuInstance,
			menuRebuild,
		),
	)

	return tpl.HTMLHandlerWithFlags(func(r *http.Request) (ib instanceRebuildContext, funcs []templating.FlagFunc, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		var instance *wisski.WissKI
		instance, err = admin.Dependencies.Instances.WissKI(r.Context(), slug)
		if err == instances.ErrWissKINotFound {
			return ib, nil, httpx.ErrNotFound
		}
		if err != nil {
			return ib, nil, err
		}

		ib.Slug = instance.Slug
		ib.System = instance.System

		// replace the menu item
		funcs = []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + instance.Slug)}),
			templating.ReplaceCrumb(menuRebuild, component.MenuItem{Title: "Rebuild", Path: template.URL("/admin/rebuild/" + instance.Slug)}),
			templating.Title(instance.Slug + " - Rebuild"),
		}

		ib.systemParams = newSystemParams()
		return
	})
}
