package auth

import (
	"context"
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"

	_ "embed"
)

//go:embed "templates/totp_enable.html"
var totpEnableStr string
var totpEnableTemplate = static.AssetsUser.MustParseShared("totp_enable.html", totpEnableStr)

func (auth *Auth) authTOTPEnable(ctx context.Context) http.Handler {
	return &httpx.Form[struct{}]{
		Fields: []httpx.Field{
			{Name: "password", Type: httpx.PasswordField, EmptyOnError: true, Label: "Current Password"},
		},
		FieldTemplate: httpx.PureCSSFieldTemplate,

		CSRF: auth.csrf.Get(nil),

		SkipForm: func(r *http.Request) (data struct{}, skip bool) {
			user, err := auth.UserOf(r)
			return struct{}{}, err == nil && user != nil && user.TOTPEnabled
		},

		RenderTemplate:        totpEnableTemplate,
		RenderTemplateContext: auth.UserFormContext,

		Validate: func(r *http.Request, values map[string]string) (struct{}, error) {
			password := values["password"]

			user, err := auth.UserOf(r)
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

func (auth *Auth) authTOTPEnroll(ctx context.Context) http.Handler {
	return &httpx.Form[struct{}]{
		Fields: []httpx.Field{
			{Name: "password", Type: httpx.PasswordField, EmptyOnError: true, Label: "Current Password"},
			{Name: "passcode", Type: httpx.TextField, EmptyOnError: true, Label: "Passcode"},
		},
		FieldTemplate: httpx.PureCSSFieldTemplate,

		CSRF: auth.csrf.Get(nil),

		SkipForm: func(r *http.Request) (data struct{}, skip bool) {
			user, _ := auth.UserOf(r)
			return struct{}{}, user != nil && user.TOTPEnabled
		},
		RenderForm: func(context httpx.FormContext, w http.ResponseWriter, r *http.Request) {
			user, err := auth.UserOf(r)

			ctx := totpEnrollContext{
				userFormContext: userFormContext{
					FormContext: context,
				},
			}

			if err == nil && user != nil {
				ctx.userFormContext.User = &user.User
				secret, err := user.TOTP()
				if err == nil {
					img, _ := TOTPLink(secret, 500, 500)

					ctx.TOTPImage = template.URL(img)
					ctx.TOTPURL = template.URL(secret.URL())
				}
			}
			httpx.WriteHTML(ctx, nil, totpEnrollTemplate, "", w, r)
		},

		Validate: func(r *http.Request, values map[string]string) (struct{}, error) {
			password, passcode := values["password"], values["passcode"]

			user, err := auth.UserOf(r)
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
				err := user.EnableTOTP(r.Context(), passcode)
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

func (auth *Auth) authTOTPDisable(ctx context.Context) http.Handler {
	return &httpx.Form[struct{}]{
		Fields: []httpx.Field{
			{Name: "password", Type: httpx.PasswordField, EmptyOnError: true, Label: "Current Password"},
			{Name: "passcode", Type: httpx.TextField, EmptyOnError: true, Label: "Current Passcode"},
		},
		FieldTemplate: httpx.PureCSSFieldTemplate,

		CSRF: auth.csrf.Get(nil),

		SkipForm: func(r *http.Request) (data struct{}, skip bool) {
			user, _ := auth.UserOf(r)
			return struct{}{}, user != nil && !user.TOTPEnabled
		},
		RenderTemplate:        totpDisableTemplate,
		RenderTemplateContext: auth.UserFormContext,

		Validate: func(r *http.Request, values map[string]string) (struct{}, error) {
			password, passcode := values["password"], values["passcode"]

			user, err := auth.UserOf(r)
			if err != nil {
				return struct{}{}, err
			}

			{
				err := user.CheckCredentials(r.Context(), []byte(password), passcode)
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
