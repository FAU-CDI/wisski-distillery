package admin

import (
	"net/http"
	"time"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
	"github.com/gorilla/mux"
)

//go:embed "html/components.html"
var componentsTemplateString string
var componentsTemplate = static.AssetsComponentsIndex.MustParseShared(
	"components.html",
	componentsTemplateString,
)

type componentContext struct {
	Time time.Time

	Analytics lazy.PoolAnalytics
}

func (admin *Admin) components(r *http.Request) (cp componentContext, err error) {
	cp.Analytics = *admin.Analytics
	cp.Time = time.Now().UTC()

	return
}

//go:embed "html/ingredients.html"
var ingredientsTemplateString string
var ingredientsTemplate = static.AssetsInstanceComponentsIndex.MustParseShared(
	"ingredients.html",
	ingredientsTemplateString,
)

type ingredientsContext struct {
	Time time.Time

	Instance  models.Instance
	Analytics *lazy.PoolAnalytics
}

func (admin *Admin) ingredients(r *http.Request) (cp ingredientsContext, err error) {
	cp.Time = time.Now().UTC()

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
