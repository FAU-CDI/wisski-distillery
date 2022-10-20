package info

import (
	"net/http"
	"time"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
)

//go:embed "html/info_components.html"
var componentsTemplateString string
var componentsTemplate = static.AssetsComponentsIndex.MustParseShared(
	"info_components.html",
	componentsTemplateString,
)

type componentsPageContext struct {
	Time time.Time

	Analytics lazy.PoolAnalytics
}

func (info *Info) componentsPageAPI(r *http.Request) (cp componentsPageContext, err error) {
	cp.Analytics = *info.Analytics
	cp.Time = time.Now().UTC()

	return
}
