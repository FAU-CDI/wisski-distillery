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
			templating.Title(instance.Slug),

			admin.instanceTabs(slug, "overview"),
		}

		return
	})
}

// instanceTabs
func (admin *Admin) instanceTabs(slug string, active string) templating.FlagFunc {
	return func(flags templating.Flags, r *http.Request) templating.Flags {
		flags.Tabs = []component.MenuItem{
			{Title: "Overview", Path: template.URL("/admin/instance/" + slug), Active: active == "overview"},
			{Title: "Rebuild", Path: template.URL("/admin/instance/" + slug + "/rebuild"), Active: active == "rebuild"},
			{Title: "Users & Grants", Path: template.URL("/admin/instance/" + slug + "/users"), Active: active == "users"},
			{Title: "Purge", Path: template.URL("/admin/instance/" + slug + "/purge"), Active: active == "purge"},

			// TODO: These still need to be migrated to their own tabs
			// Then we also need to redo the main page
			/*
				{Title: "Status", Path: template.URL("/instance/" + slug + "/status"), Active: active == "status"},
				{Title: "Database", Path: template.URL("/instance/" + slug + "/database"), Active: active == "database"},
				{Title: "Drupal", Path: template.URL("/instance/" + slug + "/drupal"), Active: active == "drupal"},
				{Title: "Users & Grants", Path: template.URL("/instance/" + slug + "/users"), Active: active == "users"},
				{Title: "Stats", Path: template.URL("/instance/" + slug + "/stats"), Active: active == "stats"},
				{Title: "SSH", Path: template.URL("/instance/" + slug + "/ssh"), Active: active == "ssh"},
				{Title: "Snapshots", Path: template.URL("/instance/" + slug + "/snapshots"), Active: active == "snapshots"},
			*/
		}
		return flags
	}
}
