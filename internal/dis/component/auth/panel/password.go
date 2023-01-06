package panel

import (
	"context"
	"errors"
	"net/http"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
)

//go:embed "templates/password.html"
var passwordHTMLString string
var passwordTemplate = static.AssetsUser.MustParseShared("password.html", passwordHTMLString)

var (
	errPasswordsNotIdentical = errors.New("passwords are not identical")
	errPasswordIsEmpty       = errors.New("password is empty")
	errCredentialsIncorrect  = errors.New("credentials are not correct")
	errPasswordSetFailure    = errors.New("error saving new password")
	errTOTPSetFailure        = errors.New("unable to disable totp")
	errPasswordSet           = errors.New("password was updated")
)

func (panel *UserPanel) routePassword(ctx context.Context) http.Handler {
	passwordTemplate := panel.Dependencies.Custom.Template(passwordTemplate)

	return &httpx.Form[struct{}]{
		Fields: []httpx.Field{
			{Name: "old", Type: httpx.PasswordField, EmptyOnError: true, Label: "Current Password"},
			{Name: "otp", Type: httpx.TextField, EmptyOnError: true, Label: "Current Passcode (optional)"},
			{Name: "new", Type: httpx.PasswordField, EmptyOnError: true, Label: "New Password"},
			{Name: "new2", Type: httpx.PasswordField, EmptyOnError: true, Label: "New Password (again)"},
		},
		FieldTemplate: httpx.PureCSSFieldTemplate,

		RenderTemplate:        passwordTemplate,
		RenderTemplateContext: panel.UserFormContext,

		Validate: func(r *http.Request, values map[string]string) (struct{}, error) {
			old, passcode, new, new2 := values["old"], values["otp"], values["new"], values["new2"]

			if new != new2 {
				return struct{}{}, errPasswordsNotIdentical
			}

			if new == "" {
				return struct{}{}, errPasswordIsEmpty
			}

			user, err := panel.Dependencies.Auth.UserOf(r)
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
