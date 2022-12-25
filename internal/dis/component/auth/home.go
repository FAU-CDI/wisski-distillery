package auth

import (
	"context"
	_ "embed"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
)

//go:embed "templates/home.html"
var homeHTMLStr string
var homeTemplate = static.AssetsAuthHome.MustParseShared(
	"home.html",
	homeHTMLStr,
)

func (auth *Auth) authHome(ctx context.Context) http.Handler {
	return auth.Protect(&httpx.HTMLHandler[*AuthUser]{
		Handler:  auth.UserOf,
		Template: homeTemplate,
	}, nil)
}
