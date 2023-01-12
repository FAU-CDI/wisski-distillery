package admin

import (
	_ "embed"
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static/custom"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/julienschmidt/httprouter"
)

//go:embed "html/instance.html"
var instanceTemplateString string
var instanceTemplate = static.AssetsAdmin.MustParseShared(
	"instance.html",
	instanceTemplateString,
)

type instanceContext struct {
	custom.BaseContext

	Instance models.Instance
	Info     status.WissKI
}

func (admin *Admin) instance(r *http.Request) (is instanceContext, err error) {
	slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

	admin.Dependencies.Custom.Update(&is, r, custom.BaseContextGaps{
		Crumbs: []component.MenuItem{
			{Title: "Admin", Path: "/admin/"},
			{Title: "Instance", Path: template.URL("/admin/instance/" + slug)},
		},
		Actions: []component.MenuItem{
			{Title: "Grants", Path: template.URL("/admin/grants/" + slug)},
			{Title: "Ingredients", Path: template.URL("/admin/ingredients/" + slug), Priority: component.SmallButton},
		},
	})

	// find the instance itself!
	instance, err := admin.Dependencies.Instances.WissKI(r.Context(), slug)
	if err == instances.ErrWissKINotFound {
		return is, httpx.ErrNotFound
	}
	if err != nil {
		return is, err
	}
	is.Instance = instance.Instance

	// get some more info about the wisski
	is.Info, err = instance.Info().Information(r.Context(), false)
	if err != nil {
		return is, err
	}

	return
}
