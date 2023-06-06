package panel

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/rs/zerolog"
	"github.com/tkw1536/pkglib/httpx"
	"github.com/tkw1536/pkglib/httpx/field"

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

	Tokens []models.Token
}

func (panel *UserPanel) tokensRoute(ctx context.Context) http.Handler {
	tpl := tokensTemplate.Prepare(
		panel.Dependencies.Templating,
		templating.Crumbs(
			menuUser,
			menuTokens,
		),
		templating.Actions(
			menuTokensAdd,
		),
	)

	return tpl.HTMLHandler(func(r *http.Request) (tc TokenTemplateContext, err error) {
		// list the user
		user, err := panel.Dependencies.Auth.UserOfSession(r)
		if err != nil || user == nil {
			return tc, err
		}

		// get the tokens
		tc.Tokens, err = panel.Dependencies.Tokens.Tokens(r.Context(), user.User.User)
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
		user, err := panel.Dependencies.Auth.UserOfSession(r)
		if err != nil {
			logger.Err(err).Str("action", "delete token").Msg("failed to get current user")
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		token := r.PostFormValue("token")
		if token == "" {
			logger.Err(err).Str("action", "delete token").Msg("failed to get token")
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		if err := panel.Dependencies.Tokens.Remove(r.Context(), user.User.User, token); err != nil {
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
	"tokens_add.html", tokensAddHTML, httpx.FormTemplate,
	templating.Title("Add Token"),
	templating.Assets(assets.AssetsUser),
)

type addTokenResult struct {
	User        *auth.AuthUser
	Description string
	Scopes      []string
}

func (panel *UserPanel) tokensAddRoute(ctx context.Context) http.Handler {
	tpl := tokensAddTemplate.Prepare(
		panel.Dependencies.Templating,
		templating.Crumbs(
			menuUser,
			menuTokens,
			menuTokensAdd,
		),
	)

	return &httpx.Form[addTokenResult]{
		Fields: []field.Field{
			{Name: "description", Type: field.Text, Label: "Description"},
		},
		FieldTemplate: field.PureCSSFieldTemplate,

		RenderTemplate:        tpl.Template(),
		RenderTemplateContext: templating.FormTemplateContext(tpl),

		Validate: func(r *http.Request, values map[string]string) (at addTokenResult, err error) {
			at.User, err = panel.Dependencies.Auth.UserOfSession(r)
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

		RenderSuccess: func(at addTokenResult, values map[string]string, w http.ResponseWriter, r *http.Request) error {
			// add the key to the user
			_, err := panel.Dependencies.Tokens.Add(r.Context(), at.User.User.User, at.Description, at.Scopes)
			if err != nil {
				return err
			}
			if err != nil {
				return errAddToken
			}
			// everything went fine, redirect the user back to the user page!
			http.Redirect(w, r, string(menuTokens.Path), http.StatusSeeOther)
			return nil
		},
	}
}
