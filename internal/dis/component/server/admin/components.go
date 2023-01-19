package admin

import (
	"context"
	"html/template"
	"net/http"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templates"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
	"github.com/julienschmidt/httprouter"
)

//go:embed "html/components.html"
var componentsHTML []byte
var componentsTemplate = templates.Parse[componentContext]("components.html", componentsHTML, assets.AssetsAdmin)

type componentContext struct {
	templates.BaseContext

	Analytics lazy.PoolAnalytics
}

func (admin *Admin) components(ctx context.Context) http.Handler {
	tpl := componentsTemplate.Prepare(admin.Dependencies.Templating, templates.BaseContextGaps{
		Crumbs: []component.MenuItem{
			{Title: "Admin", Path: "/admin/"},
			{Title: "Components", Path: "/admin/components/"},
		},
	})

	return tpl.HTMLHandler(func(r *http.Request) (cp componentContext, err error) {
		cp.Analytics = *admin.Analytics
		return
	})
}

//go:embed "html/ingredients.html"
var ingredientsHTML []byte
var ingredientsTemplate = templates.Parse[ingredientsContext]("ingredients.html", ingredientsHTML, assets.AssetsAdmin)

type ingredientsContext struct {
	templates.BaseContext

	Instance  models.Instance
	Analytics *lazy.PoolAnalytics
}

func (admin *Admin) ingredients(ctx context.Context) http.Handler {
	tpl := ingredientsTemplate.Prepare(admin.Dependencies.Templating, templates.BaseContextGaps{
		Crumbs: []component.MenuItem{
			{Title: "Admin", Path: "/admin/"},
			{Title: "Instance", Path: "* to be updated *"},
			{Title: "Ingredients", Path: "* to be updated *"},
		},
	})

	return tpl.HTMLHandlerWithGaps(func(r *http.Request, gaps *templates.BaseContextGaps) (ic ingredientsContext, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		gaps.Crumbs[1] = component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + slug)}
		gaps.Crumbs[2] = component.MenuItem{Title: "Ingredients", Path: template.URL("/admin/instance/" + slug + "/ingredients/")}

		// find the instance itself!
		instance, err := admin.Dependencies.Instances.WissKI(r.Context(), slug)
		if err == instances.ErrWissKINotFound {
			return ic, httpx.ErrNotFound
		}
		if err != nil {
			return ic, err
		}
		ic.Instance = instance.Instance

		// and get the components
		ic.Analytics = instance.Info().Analytics

		return
	})
}
