//spellchecker:words admin
package admin

//spellchecker:words context html template http github wisski distillery internal component instances server assets templating models status julienschmidt httprouter pkglib httpx embed
import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

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

func (admin *Admin) instance(context.Context) http.Handler {
	tpl := instanceTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuInstance,
		),
	)

	return tpl.HTMLHandlerWithFlags(admin.dependencies.Handling, func(r *http.Request) (ic instanceContext, funcs []templating.FlagFunc, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		// find the instance itself!
		instance, err := admin.dependencies.Instances.WissKI(r.Context(), slug)
		if errors.Is(err, instances.ErrWissKINotFound) {
			return ic, nil, httpx.ErrNotFound
		}
		if err != nil {
			return ic, nil, fmt.Errorf("failed to get instance: %w", err)
		}
		ic.Instance = instance.Instance

		// get some more info about the wisski
		ic.Info, err = instance.Info().Information(r.Context(), true)
		if err != nil {
			return ic, nil, fmt.Errorf("failed to get information: %w", err)
		}

		escapedSlug := url.PathEscape(slug)
		funcs = []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + escapedSlug)}), // #nosec G203 -- escaped and safe
			templating.Title(instance.Slug),

			admin.instanceTabs(escapedSlug, "overview"),
		}

		return
	})
}

// #nosec G203 -- escaped and safe
func (admin *Admin) instanceTabs(slugEscaped string, active string) templating.FlagFunc {
	return func(flags templating.Flags, r *http.Request) templating.Flags {
		flags.Tabs = []component.MenuItem{
			{Title: "Overview", Path: template.URL("/admin/instance/" + slugEscaped), Active: active == "overview"},
			{Title: "Rebuild", Path: template.URL("/admin/instance/" + slugEscaped + "/rebuild"), Active: active == "rebuild"},
			{Title: "Users & Grants", Path: template.URL("/admin/instance/" + slugEscaped + "/users"), Active: active == "users"},
			{Title: "Triplestore", Path: template.URL("/admin/instance/" + slugEscaped + "/triplestore"), Active: active == "triplestore"},
			{Title: "Drupal", Path: template.URL("/admin/instance/" + slugEscaped + "/drupal"), Active: active == "drupal"},
			{Title: "WissKI Data", Path: template.URL("/admin/instance/" + slugEscaped + "/data"), Active: active == "data"},
			{Title: "WissKI Stats", Path: template.URL("/admin/instance/" + slugEscaped + "/stats"), Active: active == "stats"},
			{Title: "SSH", Path: template.URL("/admin/instance/" + slugEscaped + "/ssh"), Active: active == "ssh"},
			{Title: "Snapshots", Path: template.URL("/admin/instance/" + slugEscaped + "/snapshots"), Active: active == "snapshots"},
			{Title: "Purge", Path: template.URL("/admin/instance/" + slugEscaped + "/purge"), Active: active == "purge"},
		}
		return flags
	}
}
