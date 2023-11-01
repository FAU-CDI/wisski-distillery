package legal

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"

	_ "embed"
)

type Legal struct {
	component.Base
	dependencies struct {
		Static     *assets.Static
		Templating *templating.Templating
	}
}

var (
	_ component.Routeable = (*Legal)(nil)
)

//go:embed "legal.html"
var legalHTML []byte
var legalTemplate = templating.Parse[legalContext](
	"legal.html", legalHTML, nil,

	templating.Title("Legal"),
	templating.Assets(assets.AssetsDefault),
)

type legalContext struct {
	templating.RuntimeFlags

	LegalNotices string

	CSRFCookie       string
	SessionCookie    string
	AssetsDisclaimer string
}

func (legal *Legal) Routes() component.Routes {
	return component.Routes{
		Prefix: "/legal/",
		Exact:  true,

		CSRF: false,
	}
}

var (
	menuLegal = component.MenuItem{Title: "Legal", Path: "/legal/"}
)

func (legal *Legal) HandleRoute(ctx context.Context, route string) (http.Handler, error) {
	tpl := legalTemplate.Prepare(
		legal.dependencies.Templating,
		templating.Crumbs(
			menuLegal,
		),
	)

	return tpl.HTMLHandler(func(r *http.Request) (lc legalContext, err error) {
		lc.LegalNotices = cli.LegalNotices

		lc.CSRFCookie = server.CSRFCookie
		lc.SessionCookie = server.SessionCookie
		lc.AssetsDisclaimer = assets.Disclaimer

		return
	}), nil
}
