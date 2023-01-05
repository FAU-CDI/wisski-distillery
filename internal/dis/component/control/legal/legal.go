package legal

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"

	_ "embed"
)

type Legal struct {
	component.Base
}

var (
	_ component.Routeable = (*Legal)(nil)
)

//go:embed "legal.html"
var legalTemplateString string
var legalTemplate = static.AssetsLegal.MustParseShared("legal.html", legalTemplateString)

func (legal *Legal) Routes() component.Routes {
	return component.Routes{
		Paths: []string{"/legal/"},
		CSRF:  false,
	}
}

func (legal *Legal) HandleRoute(ctx context.Context, route string) (http.Handler, error) {
	return httpx.HTMLHandler[legalContext]{
		Handler:  legal.context,
		Template: legalTemplate,
	}, nil
}

type legalContext struct {
	LegalNotices string

	CSRFCookie       string
	SessionCookie    string
	AssetsDisclaimer string
}

func (legal *Legal) context(r *http.Request) (legalContext, error) {
	return legalContext{
		LegalNotices: cli.LegalNotices,

		CSRFCookie:       control.CSRFCookie,
		SessionCookie:    control.SessionCookie,
		AssetsDisclaimer: static.AssetsDisclaimer,
	}, nil
}
