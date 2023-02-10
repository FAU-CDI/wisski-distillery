package home

import (
	"context"
	_ "embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/FAU-CDI/wisski-distillery/pkg/pools"
)

//go:embed "public.html"
var publicHTML []byte
var publicTemplate = templating.Parse[publicContext](
	"public.html", publicHTML, nil,

	templating.Title("WissKI Distillery"),
	templating.Assets(assets.AssetsDefault),
)

//go:embed "about.html"
var aboutHTML string
var aboutTemplate = template.Must(template.New("about.html").Parse(aboutHTML))

// aboutContext is passed to about.html
type aboutContext struct {
	Instances    []status.WissKI
	Logo         template.HTML
	SelfRedirect string
}

// publicCOntext is passed to public.html
type publicContext struct {
	templating.RuntimeFlags

	aboutContext
	About template.HTML
}

const logoHTML = template.HTML(`<img src="/logo.svg" alt="WissKI Distillery Logo" class="biglogo">`)

func (home *Home) publicHandler(ctx context.Context) http.Handler {

	tpl := publicTemplate.Prepare(
		home.Dependencies.Templating,
		templating.Crumbs(
			menuHome,
		),
	)

	about := home.Dependencies.Templating.GetCustomizable(aboutTemplate)

	return tpl.HTMLHandler(func(r *http.Request) (pc publicContext, err error) {
		// only act on the root path!
		if strings.TrimSuffix(r.URL.Path, "/") != "" {
			return pc, httpx.ErrNotFound
		}

		// get a builder
		builder := pools.GetBuilder()
		defer pools.ReleaseBuilder(builder)

		// prepare about
		pc.aboutContext.Logo = logoHTML
		pc.aboutContext.Instances = home.homeInstances.Get(nil)
		pc.aboutContext.SelfRedirect = home.Config.SelfRedirect.String()

		// render the about template

		if err := about.Execute(builder, pc.aboutContext); err != nil {
			return pc, nil
		}

		// and return about!
		pc.About = template.HTML(builder.String())

		return
	})
}
