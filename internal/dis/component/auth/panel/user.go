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
)

//go:embed "templates/user.html"
var userHTML []byte
var userTemplate = custom.Parse[userContext]("user.html", userHTML, static.AssetsUser)

type userContext struct {
	custom.BaseContext
	*auth.AuthUser

	Grants []GrantWithURL
}

type GrantWithURL struct {
	models.Grant
	URL template.URL
}

func (panel *UserPanel) routeUser(ctx context.Context) http.Handler {
	tpl := userTemplate.Prepare(panel.Dependencies.Custom, custom.BaseContextGaps{
		Crumbs: []component.MenuItem{
			{Title: "User", Path: "/user/"},
		},
		Actions: []component.MenuItem{
			{Title: "Change Password", Path: "/user/password/"},
			{Title: "*to be replaced*", Path: ""},
			{Title: "SSH Keys", Path: "/user/ssh/"},
		},
	})

	return tpl.HTMLHandlerWithGaps(func(r *http.Request, gaps *custom.BaseContextGaps) (uc userContext, err error) {
		// find the user
		uc.AuthUser, err = panel.Dependencies.Auth.UserOf(r)
		if err != nil || uc.AuthUser == nil {
			return uc, err
		}

		// build the gaps
		if uc.AuthUser.IsTOTPEnabled() {
			gaps.Actions[1] = component.MenuItem{
				Title: "Disable Passcode (TOTP)",
				Path:  "/user/totp/disable/",
			}
		} else {
			gaps.Actions[1] = component.MenuItem{
				Title: "Enable Passcode (TOTP)",
				Path:  "/user/totp/enable/",
			}
		}

		// find the grants
		grants, err := panel.Dependencies.Policy.User(r.Context(), uc.AuthUser.User.User)
		if err != nil {
			return uc, err
		}

		uc.Grants = make([]GrantWithURL, len(grants))
		for i, grant := range grants {
			uc.Grants[i].Grant = grant

			url, err := panel.Dependencies.Next.Next(r.Context(), grant.Slug, "/")
			if err != nil {
				return uc, err
			}
			uc.Grants[i].URL = template.URL(url)
		}

		return uc, err
	})
}
