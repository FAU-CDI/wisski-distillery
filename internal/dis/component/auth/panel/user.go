package panel

import (
	"context"
	"html/template"
	"net/http"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static/custom"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
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

	Grants []GrantWithURL
}

type GrantWithURL struct {
	models.Grant
	URL template.URL
}

func (panel *UserPanel) routeUser(ctx context.Context) http.Handler {
	userTemplate := panel.Dependencies.Custom.Template(userTemplate)
	crumbs := []component.MenuItem{
		{Title: "User", Path: "/user/"},
	}

	return &httpx.HTMLHandler[routeUserContext]{
		Handler: func(r *http.Request) (ruc routeUserContext, err error) {
			panel.Dependencies.Custom.Update(&ruc, r, crumbs)

			// find the user
			ruc.AuthUser, err = panel.Dependencies.Auth.UserOf(r)
			if err != nil || ruc.AuthUser == nil {
				return ruc, err
			}

			// find the grants
			grants, err := panel.Dependencies.Policy.User(r.Context(), ruc.AuthUser.User.User)
			if err != nil {
				return ruc, err
			}

			ruc.Grants = make([]GrantWithURL, len(grants))
			for i, grant := range grants {
				ruc.Grants[i].Grant = grant

				url, err := panel.Dependencies.Next.Next(r.Context(), grant.Slug, "/")
				if err != nil {
					return ruc, err
				}
				ruc.Grants[i].URL = template.URL(url)
			}

			return ruc, err
		},
		Template: userTemplate,
	}
}
