//spellchecker:words panel
package panel

//spellchecker:words context errors html template http github wisski distillery internal component auth server assets templating models wdlog pkglib httpx form field embed
import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"go.tkw01536.de/pkglib/httpx"
	"go.tkw01536.de/pkglib/httpx/form"
	"go.tkw01536.de/pkglib/httpx/form/field"

	_ "embed"
)

//go:embed "templates/tokens.html"
var tokensHTML []byte
var tokensTemplate = templating.Parse[TokenTemplateContext](
	"tokens.html", tokensHTML, nil,

	templating.Title("Tokens"),
	templating.Assets(assets.AssetsUser),
)

type TokenTemplateContext struct {
	templating.RuntimeFlags

	Domain template.URL // server base URL
	Tokens []models.Token
}

var errNoUserInSession = errors.New("no user in session")

func (panel *UserPanel) tokensRoute(context.Context) http.Handler {
	tpl := tokensTemplate.Prepare(
		panel.dependencies.Templating,
		templating.Crumbs(
			menuUser,
			menuTokens,
		),
		templating.Actions(
			menuTokensAdd,
		),
	)

	return tpl.HTMLHandler(panel.dependencies.Handling, func(r *http.Request) (tc TokenTemplateContext, err error) {
		// list the user
		user, err := panel.dependencies.Auth.UserOfSession(r)
		if err != nil {
			return tc, fmt.Errorf("failed to get user of session: %w", err)
		}
		if user == nil {
			return tc, errNoUserInSession
		}

		tc.Domain = template.URL(component.GetStill(panel).Config.HTTP.JoinPath().String()) // #nosec G203 -- assumed to be safe

		// get the tokens
		tc.Tokens, err = panel.dependencies.Tokens.Tokens(r.Context(), user.User.User)
		if err != nil {
			return tc, fmt.Errorf("failed to get token: %w", err)
		}
		return tc, nil
	})
}

func (panel *UserPanel) tokensDeleteRoute(ctx context.Context) http.Handler {
	logger := wdlog.Of(ctx)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			logger.Error(
				"failed to parse form",
				"error", err,
				"action", "delete token",
			)
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}
		user, err := panel.dependencies.Auth.UserOfSession(r)
		if err != nil {
			logger.Error(
				"failed to get current user",
				"error", err,
				"action", "delete token",
			)
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		id := r.PostFormValue("id")
		if id == "" {
			logger.Error(
				"failed to get token",
				"error", err,
				"action", "delete token",
			)
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		if err := panel.dependencies.Tokens.Remove(r.Context(), user.User.User, id); err != nil {
			logger.Error(
				"failed to delete token",
				"error", err,
				"action", "delete token",
			)
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, string(menuTokens.Path), http.StatusSeeOther)
	})
}

//go:embed "templates/tokens_add.html"
var tokensAddHTML []byte
var tokensAddTemplate = templating.ParseForm(
	"tokens_add.html", tokensAddHTML, form.FormTemplate,
	templating.Title("Add Token"),
	templating.Assets(assets.AssetsUser),
)

type addTokenResult struct {
	User        *auth.AuthUser
	Description string
	Scopes      []string
}

//go:embed "templates/token_created.html"
var tokenCreatedHTML []byte
var tokenCreateTemplate = templating.Parse[TokenCreateContext](
	"token_created.html", tokenCreatedHTML, form.FormTemplate,
	templating.Title("Add Token"),
	templating.Assets(assets.AssetsUser),
)

type TokenCreateContext struct {
	templating.RuntimeFlags

	Domain template.URL // server base URL
	Token  *models.Token
}

func (panel *UserPanel) tokensAddRoute(context.Context) http.Handler {
	tplForm := tokensAddTemplate.Prepare(
		panel.dependencies.Templating,
		templating.Crumbs(
			menuUser,
			menuTokens,
			menuTokensAdd,
		),
	)

	tplDone := tokenCreateTemplate.Prepare(
		panel.dependencies.Templating,
		templating.Crumbs(
			menuUser,
			menuTokens,
			menuTokensAdd,
		),
	)

	return &form.Form[addTokenResult]{
		Fields: []field.Field{
			{Name: "description", Type: field.Text, Label: "Description"},
		},
		FieldTemplate: assets.PureCSSFieldTemplate,

		Template:         tplForm.Template(),
		TemplateContext:  templating.FormTemplateContext(tplForm),
		LogTemplateError: tplForm.LogTemplateError,

		Validate: func(r *http.Request, values map[string]string) (at addTokenResult, err error) {
			at.User, err = panel.dependencies.Auth.UserOfSession(r)
			if err != nil || at.User == nil {
				return at, errInvalidUser
			}

			at.Description = values["description"]
			if at.Description == "" {
				at.Description = "API Key"
			}

			at.Scopes = nil

			return at, nil
		},

		Success: func(at addTokenResult, values map[string]string, w http.ResponseWriter, r *http.Request) error {
			// add the key to the user
			tok, err := panel.dependencies.Tokens.Add(r.Context(), at.User.User.User, at.Description, at.Scopes)
			if err != nil {
				return errAddToken
			}

			// render the created context
			return panel.dependencies.Handling.WriteHTML(
				tplDone.Context(r, TokenCreateContext{
					Domain: template.URL(component.GetStill(panel).Config.HTTP.JoinPath().String()), // #nosec G203 -- assumed to be safe
					Token:  tok,
				}),
				nil,
				tplDone.Template(),
				w,
				r,
			)
		},
	}
}
