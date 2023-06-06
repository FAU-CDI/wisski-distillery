package panel

import (
	"context"
	"html/template"
	"net/http"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

//go:embed "templates/user.html"
var userHTML []byte
var userTemplate = templating.Parse[userContext](
	"user.html", userHTML, nil,

	templating.Assets(assets.AssetsUser),
)

type userContext struct {
	templating.RuntimeFlags
	*auth.AuthUser

	ShowAdminURLs bool
	Grants        []GrantWithURL
}

type GrantWithURL struct {
	models.Grant
	URL template.URL
}

func (g GrantWithURL) AdminURL() template.URL {
	return template.URL("/admin/instance/" + g.Slug)
}

func (panel *UserPanel) routeUser(ctx context.Context) http.Handler {

	tpl := userTemplate.Prepare(
		panel.Dependencies.Templating,
		templating.Crumbs(
			menuUser,
		),
		templating.Actions(
			menuChangePassword,
			menuTOTPAction,
			menuSSH,
			menuTokens,
		),
	)

	return tpl.HTMLHandlerWithFlags(func(r *http.Request) (uc userContext, funcs []templating.FlagFunc, err error) {
		// find the user
		uc.AuthUser, err = panel.Dependencies.Auth.UserOfSession(r)
		if err != nil || uc.AuthUser == nil {
			return uc, nil, err
		}

		uc.ShowAdminURLs = panel.Dependencies.Auth.CheckScope("", scopes.ScopeAdminLoggedIn, r) == nil

		// replace the totp action in the menu
		var totpAction component.MenuItem
		if uc.AuthUser.IsTOTPEnabled() {
			totpAction = menuTOTPDisable
		} else {
			totpAction = menuTOTPEnable
		}
		funcs = []templating.FlagFunc{
			templating.ReplaceAction(menuTOTPAction, totpAction),
			templating.Title(uc.AuthUser.User.User),
		}

		// find the grants
		grants, err := panel.Dependencies.Policy.User(r.Context(), uc.AuthUser.User.User)
		if err != nil {
			return uc, nil, err
		}

		uc.Grants = make([]GrantWithURL, len(grants))
		for i, grant := range grants {
			uc.Grants[i].Grant = grant

			url, err := panel.Dependencies.Next.Next(r.Context(), grant.Slug, "/")
			if err != nil {
				return uc, nil, err
			}
			uc.Grants[i].URL = template.URL(url)
		}

		return uc, funcs, err
	})
}
