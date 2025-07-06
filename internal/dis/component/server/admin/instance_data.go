//spellchecker:words admin
package admin

//spellchecker:words context embed html template http github wisski distillery internal component server assets templating pkglib errorsx httpx julienschmidt httprouter
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
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/httpx"

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
		defer errorsx.Close(server, &err, "server")

		ctx.Pathbuilders, err = ctx.Instance.Pathbuilder().GetAll(r.Context(), server)
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to get pathbuilders: %w", err)
		}

		prefixes := ctx.Instance.Prefixes()
		ctx.NoPrefixes = prefixes.NoPrefix()
		ctx.Prefixes, err = prefixes.All(r.Context(), server)
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to get prefixes: %w", err)
		}

		escapedSlug := url.PathEscape(ctx.Instance.Slug)
		return ctx, []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + escapedSlug)}),  // #nosec G203 -- escaped and safe
			templating.ReplaceCrumb(menuData, component.MenuItem{Title: "SSH", Path: template.URL("/admin/instance/" + escapedSlug + "/data")}), // #nosec G203 -- escaped and safe
			templating.Title(ctx.Instance.Slug + " - Data"),
			admin.instanceTabs(escapedSlug, "data"),
		}, nil
	})
}
