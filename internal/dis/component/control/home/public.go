package home

import (
	"context"
	_ "embed"
	"net/http"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static/custom"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
)

//go:embed "public.html"
var publicHTML []byte
var publicTemplate = custom.Parse[publicContext]("public.html", publicHTML, static.AssetsDefault)

type publicContext struct {
	custom.BaseContext

	Instances    []status.WissKI
	SelfRedirect string
}

func (home *Home) publicHandler(ctx context.Context) http.Handler {
	tpl := publicTemplate.Prepare(home.Dependencies.Custom, custom.BaseContextGaps{
		Crumbs: []component.MenuItem{
			{Title: "WissKI Distillery", Path: "/"},
		},
	})
	return tpl.HTMLHandler(func(r *http.Request) (pc publicContext, err error) {
		// only act on the root path!
		if strings.TrimSuffix(r.URL.Path, "/") != "" {
			return pc, httpx.ErrNotFound
		}

		pc.Instances = home.homeInstances.Get(nil)
		pc.SelfRedirect = home.Config.SelfRedirect.String()

		return
	})
}
