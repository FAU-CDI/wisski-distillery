package panel

import (
	"context"
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"

	_ "embed"
)

//go:embed "templates/totp_enable.html"
var totpEnableStr string
var totpEnableTemplate = static.AssetsUser.MustParseShared("totp_enable.html", totpEnableStr)

func (panel *UserPanel) routeTOTPEnable(ctx context.Context) http.Handler {
	totpEnableTemplate := panel.Dependencies.Custom.Template(totpEnableTemplate)

	return &httpx.Form[struct{}]{
		Fields: []httpx.Field{
			{Name: "password", Type: httpx.PasswordField, EmptyOnError: true, Label: "Current Password"},
		},
		FieldTemplate: httpx.PureCSSFieldTemplate,

		SkipForm: func(r *http.Request) (data struct{}, skip bool) {
			user, err := panel.Dependencies.Auth.UserOf(r)
			return struct{}{}, err == nil && user != nil && user.IsTOTPEnabled()
		},

		RenderTemplate:        totpEnableTemplate,
		RenderTemplateContext: panel.UserFormContext,

		Validate: func(r *http.Request, values map[string]string) (struct{}, error) {
			password := values["password"]

			user, err := panel.Dependencies.Auth.UserOf(r)
			if err != nil {
				return struct{}{}, err
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

		RenderSuccess: func(_ struct{}, values map[string]string, w http.ResponseWriter, r *http.Request) error {
			http.Redirect(w, r, "/user/totp/enroll", http.StatusSeeOther)
			return nil
		},
	}
}

//go:embed "templates/totp_enroll.html"
var totpEnrollStr string
var totpEnrollTemplate = static.AssetsUser.MustParseShared("totp_enroll.html", totpEnrollStr)

type totpEnrollContext struct {
	userFormContext
	TOTPImage template.URL
	TOTPURL   template.URL
}

func (panel *UserPanel) routeTOTPEnroll(ctx context.Context) http.Handler {
	totpEnrollTemplate := panel.Dependencies.Custom.Template(totpEnrollTemplate)

	return &httpx.Form[struct{}]{
		Fields: []httpx.Field{
			{Name: "password", Type: httpx.PasswordField, EmptyOnError: true, Label: "Current Password"},
			{Name: "otp", Type: httpx.TextField, EmptyOnError: true, Label: "Passcode"},
		},
		FieldTemplate: httpx.PureCSSFieldTemplate,

		SkipForm: func(r *http.Request) (data struct{}, skip bool) {
			user, err := panel.Dependencies.Auth.UserOf(r)
			return struct{}{}, err == nil && user != nil && user.IsTOTPEnabled()
		},
		RenderForm: func(context httpx.FormContext, w http.ResponseWriter, r *http.Request) {
			// TODO: Do we want to reuse the same function here?

			user, err := panel.Dependencies.Auth.UserOf(r)

			ctx := totpEnrollContext{
				userFormContext: userFormContext{
					FormContext: context,
				},
			}
			panel.Dependencies.Custom.Update(&ctx.userFormContext)

			if err == nil && user != nil {
				ctx.userFormContext.User = &user.User
				secret, err := user.TOTP()
				if err == nil {
					img, _ := auth.TOTPLink(secret, 500, 500)

					ctx.TOTPImage = template.URL(img)
					ctx.TOTPURL = template.URL(secret.URL())
				}
			}
			httpx.WriteHTML(ctx, nil, totpEnrollTemplate, "", w, r)
		},

		Validate: func(r *http.Request, values map[string]string) (struct{}, error) {
			password, otp := values["password"], values["otp"]

			user, err := panel.Dependencies.Auth.UserOf(r)
			if err != nil {
				return struct{}{}, err
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

		RenderSuccess: func(_ struct{}, values map[string]string, w http.ResponseWriter, r *http.Request) error {
			http.Redirect(w, r, "/user/", http.StatusSeeOther)
			return nil
		},
	}
}

//go:embed "templates/totp_disable.html"
var totpDisableStr string
var totpDisableTemplate = static.AssetsUser.MustParseShared("totp_disable.html", totpDisableStr)

func (panel *UserPanel) routeTOTPDisable(ctx context.Context) http.Handler {
	totpDisableTemplate := panel.Dependencies.Custom.Template(totpDisableTemplate)

	return &httpx.Form[struct{}]{
		Fields: []httpx.Field{
			{Name: "password", Type: httpx.PasswordField, EmptyOnError: true, Label: "Current Password"},
			{Name: "otp", Type: httpx.TextField, EmptyOnError: true, Label: "Current Passcode"},
		},
		FieldTemplate: httpx.PureCSSFieldTemplate,

		SkipForm: func(r *http.Request) (data struct{}, skip bool) {
			user, err := panel.Dependencies.Auth.UserOf(r)
			return struct{}{}, err == nil && user != nil && !user.IsTOTPEnabled()
		},
		RenderTemplate:        totpDisableTemplate,
		RenderTemplateContext: panel.UserFormContext,

		Validate: func(r *http.Request, values map[string]string) (struct{}, error) {
			password, otp := values["password"], values["otp"]

			user, err := panel.Dependencies.Auth.UserOf(r)
			if err != nil {
				return struct{}{}, err
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
					return struct{}{}, errTOTPSetFailure
				}
			}

			return struct{}{}, nil
		},

		RenderSuccess: func(_ struct{}, values map[string]string, w http.ResponseWriter, r *http.Request) error {
			http.Redirect(w, r, "/user/", http.StatusSeeOther)
			return nil
		},
	}
}
