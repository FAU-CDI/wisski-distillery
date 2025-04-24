//spellchecker:words panel
package panel

//spellchecker:words context html template http github wisski distillery internal component auth server assets templating pkglib httpx form field embed
import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/tkw1536/pkglib/httpx/form"
	"github.com/tkw1536/pkglib/httpx/form/field"

	_ "embed"
)

//go:embed "templates/totp_enable.html"
var totpEnableHTML []byte
var totpEnable = templating.Parse[userFormContext](
	"totp_enable.html", totpEnableHTML, form.FormTemplate,

	templating.Title("Enable TOTP"),
	templating.Assets(assets.AssetsUser),
)

func (panel *UserPanel) routeTOTPEnable(context.Context) http.Handler {
	tpl := totpEnable.Prepare(panel.dependencies.Templating)

	return &form.Form[struct{}]{
		Fields: []field.Field{
			{Name: "password", Type: field.Password, Autocomplete: field.CurrentPassword, EmptyOnError: true, Label: "Current Password"},
		},
		FieldTemplate: assets.PureCSSFieldTemplate,

		Skip: func(r *http.Request) (data struct{}, skip bool) {
			user, err := panel.dependencies.Auth.UserOfSession(r)
			return struct{}{}, err == nil && user != nil && user.IsTOTPEnabled()
		},

		Template:         tpl.Template(),
		TemplateContext:  panel.UserFormContext(tpl, menuTOTPEnable),
		LogTemplateError: tpl.LogTemplateError,

		Validate: func(r *http.Request, values map[string]string) (struct{}, error) {
			password := values["password"]

			user, err := panel.dependencies.Auth.UserOfSession(r)
			if err != nil {
				return struct{}{}, fmt.Errorf("failed to get user of session: %w", err)
			}

			{
				err := user.CheckPassword(r.Context(), []byte(password))
				if err != nil {
					return struct{}{}, errCredentialsIncorrect
				}
			}
			{
				_, err := user.NewTOTP(r.Context())
				if err != nil {
					return struct{}{}, errTOTPSetFailure
				}
			}

			return struct{}{}, nil
		},

		Success: func(_ struct{}, values map[string]string, w http.ResponseWriter, r *http.Request) error {
			http.Redirect(w, r, "/user/totp/enroll", http.StatusSeeOther)
			return nil
		},
	}
}

//go:embed "templates/totp_enroll.html"
var totpEnrollHTML []byte
var totpEnrollTemplate = templating.Parse[totpEnrollContext](
	"totp_enroll.html", totpEnrollHTML, form.FormTemplate,

	templating.Title("Enable TOTP"),
	templating.Assets(assets.AssetsUser),
)

//nolint:errname
type totpEnrollContext struct {
	userFormContext

	TOTPSecret string
	TOTPImage  template.URL
	TOTPURL    template.URL
}

func (panel *UserPanel) routeTOTPEnroll(context.Context) http.Handler {
	tpl := totpEnrollTemplate.Prepare(
		panel.dependencies.Templating,
		templating.Crumbs(
			menuUser,
			menuTOTPEnable,
		),
	)

	return &form.Form[struct{}]{
		Fields: []field.Field{
			{Name: "password", Type: field.Password, Autocomplete: field.CurrentPassword, EmptyOnError: true, Label: "Current Password"},
			{Name: "otp", Type: field.Text, Autocomplete: field.OneTimeCode, EmptyOnError: true, Label: "Passcode"},
		},
		FieldTemplate: assets.PureCSSFieldTemplate,

		Skip: func(r *http.Request) (data struct{}, skip bool) {
			user, err := panel.dependencies.Auth.UserOfSession(r)
			return struct{}{}, err == nil && user != nil && user.IsTOTPEnabled()
		},

		Template: tpl.Template(),
		TemplateContext: func(context form.FormContext, r *http.Request) any {
			user, err := panel.dependencies.Auth.UserOfSession(r)

			ctx := totpEnrollContext{
				userFormContext: userFormContext{
					FormContext: context,
				},
			}

			if err == nil && user != nil {
				ctx.User = &user.User
				secret, err := user.TOTP()
				if err == nil {
					img, _ := auth.TOTPLink(secret, 500, 500)

					ctx.TOTPSecret = secret.Secret()
					ctx.TOTPImage = template.URL(img)        // #nosec G203 -- this is safe
					ctx.TOTPURL = template.URL(secret.URL()) // #nosec G203 -- this is safe
				}
			}

			return tpl.Context(r, ctx)
		},
		LogTemplateError: tpl.LogTemplateError,

		Validate: func(r *http.Request, values map[string]string) (struct{}, error) {
			password, otp := values["password"], values["otp"]

			user, err := panel.dependencies.Auth.UserOfSession(r)
			if err != nil {
				return struct{}{}, fmt.Errorf("failed to get user of session: %w", err)
			}

			{
				err := user.CheckPassword(r.Context(), []byte(password))
				if err != nil {
					return struct{}{}, errCredentialsIncorrect
				}
			}
			{
				err := user.EnableTOTP(r.Context(), otp)
				if err != nil {
					return struct{}{}, errTOTPSetFailure
				}
			}

			return struct{}{}, nil
		},

		Success: func(_ struct{}, values map[string]string, w http.ResponseWriter, r *http.Request) error {
			http.Redirect(w, r, "/user/", http.StatusSeeOther)
			return nil
		},
	}
}

//go:embed "templates/totp_disable.html"
var totpDisableHTML []byte
var totpDisableTemplate = templating.Parse[userFormContext](
	"totp_disable.html", totpDisableHTML, form.FormTemplate,

	templating.Title("Disable TOTP"),
	templating.Assets(assets.AssetsUser),
)

func (panel *UserPanel) routeTOTPDisable(context.Context) http.Handler {
	tpl := totpDisableTemplate.Prepare(panel.dependencies.Templating)

	return &form.Form[struct{}]{
		Fields: []field.Field{
			{Name: "password", Type: field.Password, Autocomplete: field.CurrentPassword, EmptyOnError: true, Label: "Current Password"},
			{Name: "otp", Type: field.Text, Autocomplete: field.OneTimeCode, EmptyOnError: true, Label: "Current Passcode"},
		},
		FieldTemplate: assets.PureCSSFieldTemplate,

		Skip: func(r *http.Request) (data struct{}, skip bool) {
			user, err := panel.dependencies.Auth.UserOfSession(r)
			return struct{}{}, err == nil && user != nil && !user.IsTOTPEnabled()
		},

		Template:         tpl.Template(),
		TemplateContext:  panel.UserFormContext(tpl, menuTOTPDisable),
		LogTemplateError: tpl.LogTemplateError,

		Validate: func(r *http.Request, values map[string]string) (struct{}, error) {
			password, otp := values["password"], values["otp"]

			user, err := panel.dependencies.Auth.UserOfSession(r)
			if err != nil {
				return struct{}{}, fmt.Errorf("failed to get user of session: %w", err)
			}

			{
				err := user.CheckCredentials(r.Context(), []byte(password), otp)
				if err != nil {
					return struct{}{}, errCredentialsIncorrect
				}
			}
			{
				err := user.DisableTOTP(r.Context())
				if err != nil {
					return struct{}{}, errTOTPUnsetFailure
				}
			}

			return struct{}{}, nil
		},

		Success: func(_ struct{}, values map[string]string, w http.ResponseWriter, r *http.Request) error {
			http.Redirect(w, r, "/user/", http.StatusSeeOther)
			return nil
		},
	}
}
