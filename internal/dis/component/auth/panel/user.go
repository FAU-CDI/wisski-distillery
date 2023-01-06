package panel

import (
	"context"
	"net/http"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static/custom"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
)

//go:embed "templates/user.html"
var userHTMLStr string
var userTemplate = static.AssetsUser.MustParseShared(
	"user.html",
	userHTMLStr,
)

type routeUserContext struct {
	custom.BaseContext
	*auth.AuthUser
}

func (panel *UserPanel) routeUser(ctx context.Context) http.Handler {
	userTemplate := panel.Dependencies.Custom.Template(userTemplate)
	return &httpx.HTMLHandler[routeUserContext]{
		Handler: func(r *http.Request) (ruc routeUserContext, err error) {
			panel.Dependencies.Custom.Update(&ruc, r)
			ruc.AuthUser, err = panel.Dependencies.Auth.UserOf(r)
			return ruc, err
		},
		Template: userTemplate,
	}
}
