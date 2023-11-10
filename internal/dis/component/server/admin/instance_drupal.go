package admin

import (
	"context"
	_ "embed"
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/pkglib/httpx"

	"github.com/julienschmidt/httprouter"
)

//go:embed "html/instance_drupal.html"
var instanceDrupalHTML []byte
var instanceDrupalTemplate = templating.Parse[instanceDrupalContext](
	"instance_drupal.html", instanceDrupalHTML, nil,

	templating.Assets(assets.AssetsAdmin),
)

type instanceDrupalContext struct {
	templating.RuntimeFlags

	Instance *wisski.WissKI

	DrupalVersion string
	DefaultTheme  string

	Requirements []status.Requirement
}

func (admin *Admin) instanceDrupal(ctx context.Context) http.Handler {
	tpl := instanceDrupalTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuInstance,
			menuDrupal,
		),
	)

	return tpl.HTMLHandlerWithFlags(func(r *http.Request) (ctx instanceDrupalContext, funcs []templating.FlagFunc, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		// setup the context with just the instance
		ctx.Instance, err = admin.dependencies.Instances.WissKI(r.Context(), slug)
		if err != nil {
			return ctx, nil, httpx.ErrNotFound
		}

		server := ctx.Instance.PHP().NewServer()
		defer server.Close()

		// get the requirements
		ctx.Requirements, err = ctx.Instance.Requirements().Get(r.Context(), server)
		if err != nil {
			return ctx, nil, httpx.ErrInternalServerError
		}

		// get the drupal version
		ctx.DrupalVersion, err = ctx.Instance.Version().Get(r.Context(), server)
		if err != nil {
			return ctx, nil, httpx.ErrInternalServerError
		}

		// get the default theme
		ctx.DefaultTheme, err = ctx.Instance.Theme().Get(r.Context(), server)
		if err != nil {
			return ctx, nil, httpx.ErrInternalServerError
		}

		return ctx, []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + ctx.Instance.Slug)}),
			templating.ReplaceCrumb(menuDrupal, component.MenuItem{Title: "Drupal Status", Path: template.URL("/admin/instance/" + ctx.Instance.Slug + "/drupal")}),
			templating.Title(ctx.Instance.Slug + " - Drupal Status"),
			admin.instanceTabs(slug, "drupal"),
		}, nil
	})
}
