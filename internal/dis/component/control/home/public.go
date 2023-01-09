package home

import (
	"context"
	_ "embed"
	"net/http"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static/custom"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
)

//go:embed "public.html"
var publicHTMLStr string
var publicTemplate = static.AssetsHome.MustParseShared("public.html", publicHTMLStr)

type publicContext struct {
	custom.BaseContext

	Instances    []status.WissKI
	SelfRedirect string
}

func (home *Home) publicHandler(ctx context.Context) http.Handler {
	return httpx.HTMLHandler[publicContext]{
		Handler: func(r *http.Request) (pc publicContext, err error) {
			// only act on the root path!
			if strings.TrimSuffix(r.URL.Path, "/") != "" {
				return pc, httpx.ErrNotFound
			}

			home.Dependencies.Custom.Update(&pc, r)

			pc.Instances = home.homeInstances.Get(nil)
			pc.SelfRedirect = home.Config.SelfRedirect.String()

			return
		},
		Template: home.Dependencies.Custom.Template(publicTemplate),
	}
}
