package panel

import (
	"context"
	"html/template"
	"net/http"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
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

	Grants []GrantWithURL
}

type GrantWithURL struct {
	models.Grant
	URL template.URL
}

var (
	totpActionItem = component.DummyMenuItem()
)

func (panel *UserPanel) routeUser(ctx context.Context) http.Handler {

	tpl := userTemplate.Prepare(
		panel.Dependencies.Templating,
		templating.Crumbs(
			component.MenuItem{Title: "User", Path: "/user/"},
		),
		templating.Actions(
			component.MenuItem{Title: "Change Password", Path: "/user/password/"},
			totpActionItem,
			component.MenuItem{Title: "SSH Keys", Path: "/user/ssh/"},
		),
	)

	return tpl.HTMLHandlerWithFlags(func(r *http.Request) (uc userContext, funcs []templating.FlagFunc, err error) {
		// find the user
		uc.AuthUser, err = panel.Dependencies.Auth.UserOf(r)
		if err != nil || uc.AuthUser == nil {
			return uc, nil, err
		}

		// replace the totp action in the menu
		var totpAction component.MenuItem
		if uc.AuthUser.IsTOTPEnabled() {
			totpAction = component.MenuItem{
				Title: "Disable Passcode (TOTP)",
				Path:  "/user/totp/disable/",
			}
		} else {
			totpAction = component.MenuItem{
				Title: "Enable Passcode (TOTP)",
				Path:  "/user/totp/enable/",
			}
		}
		funcs = []templating.FlagFunc{
			templating.ReplaceAction(totpActionItem, totpAction),
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
