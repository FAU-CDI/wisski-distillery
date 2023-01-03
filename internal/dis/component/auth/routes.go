package auth

import (
	"context"
	_ "embed"
	"errors"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
)

//go:embed "templates/home.html"
var homeHTMLStr string
var homeTemplate = static.AssetsAuthHome.MustParseShared(
	"home.html",
	homeHTMLStr,
)

func (auth *Auth) authHome(ctx context.Context) http.Handler {
	return auth.Protect(&httpx.HTMLHandler[*AuthUser]{
		Handler:  auth.UserOf,
		Template: homeTemplate,
	}, nil)
}

//go:embed "templates/password.html"
var passwordHTMLString string
var passwordTemplate = static.AssetsAuthLogin.MustParseShared("password.html", passwordHTMLString)

var (
	errPasswordsNotIdentical = errors.New("passwords are not identical")
	errPasswordIsEmpty       = errors.New("password is empty")
	errCredentialsIncorrect  = errors.New("credentials are not correct")
	errPasswordSetFailure    = errors.New("error saving new password")
	errTOTPSetFailure        = errors.New("unable to disable totp")
	errPasswordSet           = errors.New("password was updated")
)

func (auth *Auth) authPassword(ctx context.Context) http.Handler {
	return &httpx.Form[struct{}]{
		Fields: []httpx.Field{
			{Name: "old", Type: httpx.PasswordField, EmptyOnError: true, Label: "Current Password"},
			{Name: "passcode", Type: httpx.TextField, EmptyOnError: true, Label: "Current Passcode (optional)"},
			{Name: "new", Type: httpx.PasswordField, EmptyOnError: true, Label: "New Password"},
			{Name: "new2", Type: httpx.PasswordField, EmptyOnError: true, Label: "New Password (again)"},
		},
		FieldTemplate: httpx.PureCSSFieldTemplate,

		CSRF: auth.csrf.Get(nil),

		RenderTemplate: passwordTemplate,

		Validate: func(r *http.Request, values map[string]string) (struct{}, error) {
			old, passcode, new, new2 := values["old"], values["passcode"], values["new"], values["new2"]

			if new != new2 {
				return struct{}{}, errPasswordsNotIdentical
			}

			if new == "" {
				return struct{}{}, errPasswordIsEmpty
			}

			user, err := auth.UserOf(r)
			if err != nil {
				return struct{}{}, err
			}

			{
				err := user.CheckCredentials(r.Context(), []byte(old), passcode)
				if err != nil {
					return struct{}{}, errCredentialsIncorrect
				}
			}
			{
				err := user.SetPassword(r.Context(), []byte(new))
				if err != nil {
					return struct{}{}, errPasswordSetFailure
				}
			}

			return struct{}{}, nil
		},

		RenderSuccess: func(_ struct{}, values map[string]string, w http.ResponseWriter, r *http.Request) error {
			return errPasswordSet
		},
	}
}