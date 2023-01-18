package admin

import (
	"context"
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
var instanceHTML []byte
var instanceTemplate = custom.Parse[instanceContext]("instance.html", instanceHTML, static.AssetsAdmin)

type instanceContext struct {
	custom.BaseContext

	Instance models.Instance
	Info     status.WissKI
}

func (admin *Admin) instance(ctx context.Context) http.Handler {
	tpl := instanceTemplate.Prepare(admin.Dependencies.Custom, custom.BaseContextGaps{
		Crumbs: []component.MenuItem{
			{Title: "Admin", Path: "/admin/"},
			{Title: "Instance", Path: "*to be replaced*"},
		},
		Actions: []component.MenuItem{
			{Title: "Grants", Path: "*to be replaced*"},
			{Title: "Ingredients", Path: "*to be replaced*", Priority: component.SmallButton},
		},
	})

	return tpl.HTMLHandlerWithGaps(func(r *http.Request, gaps *custom.BaseContextGaps) (ic instanceContext, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		gaps.Crumbs[1] = component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + slug)}

		gaps.Actions[0] = component.MenuItem{Title: "Grants", Path: template.URL("/admin/grants/" + slug)}
		gaps.Actions[1] = component.MenuItem{Title: "Ingredients", Path: template.URL("/admin/ingredients/" + slug), Priority: component.SmallButton}

		// find the instance itself!
		instance, err := admin.Dependencies.Instances.WissKI(r.Context(), slug)
		if err == instances.ErrWissKINotFound {
			return ic, httpx.ErrNotFound
		}
		if err != nil {
			return ic, err
		}
		ic.Instance = instance.Instance

		// get some more info about the wisski
		ic.Info, err = instance.Info().Information(r.Context(), false)
		if err != nil {
			return ic, err
		}

		return
	})
}
