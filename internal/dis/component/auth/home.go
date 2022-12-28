package auth

import (
	"context"
	_ "embed"
	"errors"
	"html/template"
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

type authpasswordContext struct {
	Message string
	Form    template.HTML
}

var (
	errPasswordsNotIdentical = errors.New("passwords are not identical")
	errPasswordIsEmpty       = errors.New("password is empty")
	errPasswordIncorrect     = errors.New("old password is not correct")
	errPasswordSetFailure    = errors.New("error saving new password")
	errPasswordSet           = errors.New("password was updated")
)

func (auth *Auth) authPassword(ctx context.Context) http.Handler {
	return &httpx.Form[struct{}]{
		Fields: []httpx.Field{
			{Name: "old", Type: httpx.PasswordField, EmptyOnError: true, Label: "Current Password"},
			{Name: "new", Type: httpx.PasswordField, EmptyOnError: true, Label: "New Password"},
			{Name: "new2", Type: httpx.PasswordField, EmptyOnError: true, Label: "New Password (again)"},
		},
		FieldTemplate: httpx.PureCSSFieldTemplate,

		CSRF: auth.csrf.Get(nil),

		RenderForm: func(template template.HTML, err error, w http.ResponseWriter, r *http.Request) {
			ctx := authpasswordContext{
				Message: "",
				Form:    template,
			}
			if err != nil {
				ctx.Message = err.Error()
			}
			httpx.WriteHTML(ctx, nil, passwordTemplate, "", w, r)
		},

		Validate: func(r *http.Request, values map[string]string) (struct{}, error) {
			old, new, new2 := values["old"], values["new"], values["new2"]

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
				err := user.CheckPassword(r.Context(), []byte(old))
				if err != nil {
					return struct{}{}, errPasswordIncorrect
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
