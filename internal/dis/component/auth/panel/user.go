package panel

import (
	"context"
	"net/http"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
)

//go:embed "templates/user.html"
var userHTMLStr string
var userTemplate = static.AssetsUser.MustParseShared(
	"user.html",
	userHTMLStr,
)

func (panel *UserPanel) routeUser(ctx context.Context) http.Handler {
	return &httpx.HTMLHandler[*auth.AuthUser]{
		Handler:  panel.Dependencies.Auth.UserOf,
		Template: userTemplate,
	}
}
