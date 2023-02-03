package admin

import (
	"context"
	"html/template"
	"net/http"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
	"github.com/julienschmidt/httprouter"
)

//go:embed "html/anal.html"
var analHTML []byte
var analTemplate = templating.Parse[analContext](
	"anal.html", analHTML, nil,

	templating.Assets(assets.AssetsAdmin),
)

type analContext struct {
	templating.RuntimeFlags

	Analytics lazy.PoolAnalytics
}

func (admin *Admin) components(ctx context.Context) http.Handler {
	tpl := analTemplate.Prepare(
		admin.Dependencies.Templating,
		templating.Crumbs(
			component.MenuItem{Title: "Admin", Path: "/admin/"},
			component.MenuItem{Title: "Instances", Path: "/admin/instance/"},
			component.MenuItem{Title: "Components", Path: "/admin/components/"},
		),
		templating.Title("Components"),
	)

	return tpl.HTMLHandler(func(r *http.Request) (ac analContext, err error) {
		ac.Analytics = *admin.Analytics
		return
	})
}

func (admin *Admin) ingredients(ctx context.Context) http.Handler {
	tpl := analTemplate.Prepare(
		admin.Dependencies.Templating,
		templating.Crumbs(
			component.MenuItem{Title: "Admin", Path: "/admin/"},
			component.MenuItem{Title: "Instances", Path: "/admin/instance/"},
			component.DummyMenuItem,
			component.DummyMenuItem,
		),
	)

	return tpl.HTMLHandlerWithFlags(func(r *http.Request) (ac analContext, funcs []templating.FlagFunc, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		// find the instance itself!
		instance, err := admin.Dependencies.Instances.WissKI(r.Context(), slug)
		if err == instances.ErrWissKINotFound {
			return ac, nil, httpx.ErrNotFound
		}
		if err != nil {
			return ac, nil, err
		}
		funcs = []templating.FlagFunc{
			templating.ReplaceCrumb(2, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + slug)}),
			templating.ReplaceCrumb(3, component.MenuItem{Title: "Ingredients", Path: template.URL("/admin/instance/" + slug + "/ingredients/")}),
			templating.Title(instance.Name() + " - Ingredients"),
		}

		// and get the components
		ac.Analytics = *instance.Info().Analytics

		return
	})
}
