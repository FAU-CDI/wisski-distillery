package panel

import (
	"context"
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/rs/zerolog"
	"github.com/tkw1536/pkglib/httpx"
	"github.com/tkw1536/pkglib/httpx/form"
	"github.com/tkw1536/pkglib/httpx/form/field"

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
		if err != nil || user == nil {
			return tc, err
		}

		tc.Domain = template.URL(component.GetStill(panel).Config.HTTP.JoinPath().String())

		// get the tokens
		tc.Tokens, err = panel.dependencies.Tokens.Tokens(r.Context(), user.User.User)
		return tc, err
	})
}

func (panel *UserPanel) tokensDeleteRoute(ctx context.Context) http.Handler {
	logger := zerolog.Ctx(ctx)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			logger.Err(err).Str("action", "delete token").Msg("failed to parse form")
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}
		user, err := panel.dependencies.Auth.UserOfSession(r)
		if err != nil {
			logger.Err(err).Str("action", "delete token").Msg("failed to get current user")
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		id := r.PostFormValue("id")
		if id == "" {
			logger.Err(err).Str("action", "delete token").Msg("failed to get token")
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		if err := panel.dependencies.Tokens.Remove(r.Context(), user.User.User, id); err != nil {
			logger.Err(err).Str("action", "delete token").Msg("failed to delete token")
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

		Template:        tplForm.Template(),
		TemplateContext: templating.FormTemplateContext(tplForm),

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
					Domain: template.URL(component.GetStill(panel).Config.HTTP.JoinPath().String()),
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
