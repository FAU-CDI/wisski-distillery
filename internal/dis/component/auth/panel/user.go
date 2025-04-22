//spellchecker:words panel
package panel

//spellchecker:words context html template http embed github wisski distillery internal component auth scopes server assets templating models
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

func (panel *UserPanel) routeUser(context.Context) http.Handler {
	actions := []component.MenuItem{
		menuChangePassword,
		menuTOTPAction,
		menuSSH,
	}
	if component.GetStill(panel).Config.HTTP.API.Value {
		actions = append(actions, menuTokens)
	}

	tpl := userTemplate.Prepare(
		panel.dependencies.Templating,
		templating.Crumbs(
			menuUser,
		),
		templating.Actions(actions...),
	)

	return tpl.HTMLHandlerWithFlags(panel.dependencies.Handling, func(r *http.Request) (uc userContext, funcs []templating.FlagFunc, err error) {
		// find the user
		uc.AuthUser, err = panel.dependencies.Auth.UserOfSession(r)
		if err != nil || uc.AuthUser == nil {
			return uc, nil, err
		}

		uc.ShowAdminURLs = panel.dependencies.Auth.CheckScope("", scopes.ScopeUserAdmin, r) == nil

		// replace the totp action in the menu
		var totpAction component.MenuItem
		if uc.IsTOTPEnabled() {
			totpAction = menuTOTPDisable
		} else {
			totpAction = menuTOTPEnable
		}
		funcs = []templating.FlagFunc{
			templating.ReplaceAction(menuTOTPAction, totpAction),
			templating.Title(uc.User.User),
		}

		// find the grants
		grants, err := panel.dependencies.Policy.User(r.Context(), uc.User.User)
		if err != nil {
			return uc, nil, err
		}

		uc.Grants = make([]GrantWithURL, len(grants))
		for i, grant := range grants {
			uc.Grants[i].Grant = grant

			url, err := panel.dependencies.Next.Next(r.Context(), grant.Slug, "/")
			if err != nil {
				return uc, nil, err
			}
			uc.Grants[i].URL = template.URL(url)
		}

		return uc, funcs, err
	})
}
