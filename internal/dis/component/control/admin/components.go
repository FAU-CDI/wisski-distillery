package admin

import (
	"net/http"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static/custom"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
	"github.com/gorilla/mux"
)

//go:embed "html/components.html"
var componentsTemplateString string
var componentsTemplate = static.AssetsAdmin.MustParseShared(
	"components.html",
	componentsTemplateString,
)

type componentContext struct {
	custom.BaseContext

	Analytics lazy.PoolAnalytics
}

func (admin *Admin) components(r *http.Request) (cp componentContext, err error) {
	admin.Dependencies.Custom.Update(&cp, r)

	cp.Analytics = *admin.Analytics
	return
}

//go:embed "html/ingredients.html"
var ingredientsTemplateString string
var ingredientsTemplate = static.AssetsAdmin.MustParseShared(
	"ingredients.html",
	ingredientsTemplateString,
)

type ingredientsContext struct {
	custom.BaseContext

	Instance  models.Instance
	Analytics *lazy.PoolAnalytics
}

func (admin *Admin) ingredients(r *http.Request) (cp ingredientsContext, err error) {
	admin.Dependencies.Custom.Update(&cp, r)

	// find the instance itself!
	instance, err := admin.Dependencies.Instances.WissKI(r.Context(), mux.Vars(r)["slug"])
	if err == instances.ErrWissKINotFound {
		return cp, httpx.ErrNotFound
	}
	if err != nil {
		return cp, err
	}
	cp.Instance = instance.Instance

	// and get the components
	cp.Analytics = instance.Info().Analytics

	return
}
