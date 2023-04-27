package home

import (
	"context"
	_ "embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/tkw1536/pkglib/httpx"
)

//go:embed "public.html"
var publicHTML []byte
var publicTemplate = templating.Parse[publicContext](
	"public.html", publicHTML, nil,

	templating.Assets(assets.AssetsDefault),
)

//go:embed "about.html"
var aboutHTML string
var aboutTemplate = template.Must(template.New("about.html").Parse(aboutHTML))

// aboutContext is passed to about.html
type aboutContext struct {
	Instances    []status.WissKI // list of WissKI Instancaes
	SignedIn     bool            // is there a signed in user?
	Logo         template.HTML
	SelfRedirect string
}

// publicCOntext is passed to public.html
type publicContext struct {
	templating.RuntimeFlags

	aboutContext

	ListEnabled bool   // is the list of instances enabled?
	ListTitle   string // what is the title of the list of instances?

	About template.HTML
}

const logoHTML = template.HTML(`<img src="/logo.svg" alt="WissKI Distillery Logo" class="biglogo">`)

func (home *Home) publicHandler(ctx context.Context) http.Handler {
	title := home.Config.Home.Title

	tpl := publicTemplate.Prepare(
		home.Dependencies.Templating,
		// set title and menu item
		templating.Title(title),
		templating.Crumbs(
			component.MenuItem{Title: title, Path: "/"},
		),
	)

	about := home.Dependencies.Templating.GetCustomizable(aboutTemplate)

	return tpl.HTMLHandler(func(r *http.Request) (pc publicContext, err error) {
		// only act on the root path!
		if strings.TrimSuffix(r.URL.Path, "/") != "" {
			return pc, httpx.ErrNotFound
		}

		// get a builder
		var builder strings.Builder

		// prepare about
		pc.aboutContext.Logo = logoHTML
		pc.aboutContext.Instances = home.homeInstances.Get(nil)
		pc.aboutContext.SelfRedirect = home.Config.Home.SelfRedirect.String()
		{
			user, _ := home.Dependencies.Auth.UserOf(r)
			pc.aboutContext.SignedIn = user != nil
		}

		// render the about template

		if err := about.Execute(&builder, pc.aboutContext); err != nil {
			return pc, nil
		}

		// and return about!
		pc.About = template.HTML(builder.String())

		// user is not signed in!

		if pc.aboutContext.SignedIn {
			pc.ListEnabled = home.Config.Home.List.Private.Value
		} else {
			pc.ListEnabled = home.Config.Home.List.Public.Value
		}

		// title of the list
		pc.ListTitle = home.Config.Home.List.Title

		return
	})
}
