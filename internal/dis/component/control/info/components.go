package info

import (
	_ "embed"
	"html/template"
	"net/http"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
)

//go:embed "html/base.html"
var baseTemplateString string
var baseTemplate = template.Must(template.New("base.html").Parse(baseTemplateString))

func base(name string) *template.Template {
	clone := template.Must(baseTemplate.Clone())
	clone.Tree.Name = name
	return clone
}

//go:embed "html/info_components.html"
var componentsTemplateString string
var componentsTemplate = static.AssetsComponentsIndex.MustParse(
	base("info_components.html"),
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
