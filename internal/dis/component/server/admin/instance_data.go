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

//go:embed "html/instance_data.html"
var instanceDataHTML []byte
var instanceDataTemplate = templating.Parse[instanceDataContext](
	"instance_data.html", instanceDataHTML, nil,

	templating.Assets(assets.AssetsAdmin),
)

type instanceDataContext struct {
	templating.RuntimeFlags

	Instance     *wisski.WissKI
	Pathbuilders map[string]string
	NoPrefixes   bool
	Prefixes     []string
}

func (admin *Admin) instanceData(context.Context) http.Handler {
	tpl := instanceDataTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuInstance,
			menuData,
		),
	)

	return tpl.HTMLHandlerWithFlags(admin.dependencies.Handling, func(r *http.Request) (ctx instanceDataContext, funcs []templating.FlagFunc, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		// setup the context with just the instance
		ctx.Instance, err = admin.dependencies.Instances.WissKI(r.Context(), slug)
		if err != nil {
			return ctx, nil, httpx.ErrNotFound
		}

		server := ctx.Instance.PHP().NewServer()
		defer server.Close()

		ctx.Pathbuilders, err = ctx.Instance.Pathbuilder().GetAll(r.Context(), server)
		if err != nil {
			return ctx, nil, err
		}

		prefixes := ctx.Instance.Prefixes()
		ctx.NoPrefixes = prefixes.NoPrefix()
		ctx.Prefixes, err = prefixes.All(r.Context(), server)
		if err != nil {
			return ctx, nil, err
		}

		return ctx, []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + ctx.Instance.Slug)}),
			templating.ReplaceCrumb(menuData, component.MenuItem{Title: "SSH", Path: template.URL("/admin/instance/" + ctx.Instance.Slug + "/data")}),
			templating.Title(ctx.Instance.Slug + " - Data"),
			admin.instanceTabs(slug, "data"),
		}, nil
	})
}
