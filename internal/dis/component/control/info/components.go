package info

import (
	"net/http"
	"strings"
	"time"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
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

func (info *Info) components(r *http.Request) (cp componentContext, err error) {
	cp.Analytics = *info.Analytics
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

func (info *Info) ingredients(r *http.Request) (cp ingredientsContext, err error) {
	cp.Time = time.Now().UTC()

	// find the slug as the last component of path!
	slug := strings.TrimSuffix(r.URL.Path, "/")
	slug = slug[strings.LastIndex(slug, "/")+1:]

	// find the instance itself!
	instance, err := info.Instances.WissKI(slug)
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
