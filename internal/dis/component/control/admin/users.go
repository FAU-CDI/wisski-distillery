package admin

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static/custom"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx/field"
	"github.com/rs/zerolog"
)

//go:embed "html/users.html"
var userTemplateString string
var userTemplate = static.AssetsAdmin.MustParseShared(
	"users.html",
	userTemplateString,
)

type userContext struct {
	custom.BaseContext

	Error string
	Users []*auth.AuthUser
}

func (admin *Admin) users(r *http.Request) (uc userContext, err error) {
	admin.Dependencies.Custom.Update(&uc, r, custom.BaseContextGaps{
		Crumbs: []component.MenuItem{
			{Title: "Admin", Path: "/admin/"},
			{Title: "Users", Path: "/admin/users/"},
		},
		Actions: []component.MenuItem{
			{Title: "Create New", Path: "/admin/users/create/"},
		},
	})

	uc.Error = r.URL.Query().Get("error")
	uc.Users, err = admin.Dependencies.Auth.Users(r.Context())
	return
}

//go:embed "html/user_create.html"
var userCreateTemplateString string
var userCreateTemplate = static.AssetsAdmin.MustParseShared(
	"user_create.html",
	userCreateTemplateString,
)

var (
	errCreateInvalidUsername = errors.New("invalid username")
	errCreateInvalidPassword = errors.New("invalid password")
)

type createUserResult struct {
	User      string
	Passsword string
	Admin     bool
}

func (admin *Admin) createUser(ctx context.Context) http.Handler {
	userCreateTemplate := admin.Dependencies.Custom.Template(userCreateTemplate)
	gaps := custom.BaseContextGaps{
		Crumbs: []component.MenuItem{
			{Title: "Admin", Path: "/admin/"},
			{Title: "Users", Path: "/admin/users"},
			{Title: "Create", Path: "/admin/users/create"},
		},
	}

	return &httpx.Form[createUserResult]{
		Fields: []field.Field{
			{Name: "username", Type: field.Text, Autocomplete: field.Username, Label: "Username"},
			{Name: "password", Type: field.Password, Autocomplete: field.NewPassword, Label: "Password"},
			{Name: "admin", Type: field.Checkbox, Label: "Distillery Administrator"},
		},
		FieldTemplate: field.PureCSSFieldTemplate,

		RenderTemplate: userCreateTemplate,
		RenderTemplateContext: func(ctx httpx.FormContext, r *http.Request) any {
			return admin.Dependencies.Custom.NewForm(ctx, r, gaps)
		},

		Validate: func(r *http.Request, values map[string]string) (cu createUserResult, err error) {
			cu.User, cu.Passsword, cu.Admin = values["username"], values["password"], values["admin"] == field.CheckboxChecked

			if cu.User == "" {
				return cu, errCreateInvalidUsername
			}
			if cu.Passsword == "" {
				return cu, errCreateInvalidPassword
			}

			// check the password policy
			err = admin.Dependencies.Auth.CheckPasswordPolicy(cu.Passsword, cu.User)
			if err != nil {
				return cu, err
			}

			return cu, nil
		},

		RenderSuccess: func(cu createUserResult, values map[string]string, w http.ResponseWriter, r *http.Request) error {
			// create the user
			user, err := admin.Dependencies.Auth.CreateUser(r.Context(), cu.User)
			if err != nil {
				return err
			}

			// disable the user and setup the admin flag
			user.SetAdmin(cu.Admin)
			if err := user.Save(r.Context()); err != nil {
				return err
			}

			// set the password!
			err = user.SetPassword(r.Context(), []byte(cu.Passsword))
			if err != nil {
				return err
			}

			// everything went fine, redirect the user back to the user page!
			http.Redirect(w, r, "/admin/users/", http.StatusSeeOther)
			return nil
		},
	}
}

var errNotCurrentUser = httpx.Response{
	Body:       []byte("attempt to modify current user"),
	StatusCode: http.StatusBadRequest,
}

func (admin *Admin) useraction(ctx context.Context, name string, action func(r *http.Request, user *auth.AuthUser) error) http.Handler {
	logger := zerolog.Ctx(ctx)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			logger.Err(err).Str("action", name).Msg("failed to parse form")
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		username := r.PostFormValue("user")
		user, err := admin.Dependencies.Auth.User(r.Context(), username)
		if err != nil {
			logger.Err(err).Str("action", name).Msg("failed to get user")
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		me, err := admin.Dependencies.Auth.UserOf(r)
		if err != nil {
			logger.Err(err).Str("action", name).Msg("failed to get current user")
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		// don't allow the current user
		if me.User.User == user.User.User {
			errNotCurrentUser.ServeHTTP(w, r)
			return
		}

		if err := action(r, user); err != nil {
			logger.Err(err).Str("action", name).Msg("failed to act on user")
			http.Redirect(w, r, "/admin/users/?error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/admin/users/", http.StatusSeeOther)
	})
}

func (admin *Admin) usersDeleteHandler(ctx context.Context) http.Handler {
	return admin.useraction(ctx, "delete user", func(r *http.Request, user *auth.AuthUser) error {
		return user.Delete(r.Context())
	})
}

func (admin *Admin) usersDisableHandler(ctx context.Context) http.Handler {
	return admin.useraction(ctx, "disable user", func(r *http.Request, user *auth.AuthUser) error {
		return user.UnsetPassword(r.Context())
	})
}

func (admin *Admin) usersDisableTOTPHandler(ctx context.Context) http.Handler {
	return admin.useraction(ctx, "disable user totp", func(r *http.Request, user *auth.AuthUser) error {
		return user.DisableTOTP(r.Context())
	})
}

func (admin *Admin) usersToggleAdmin(ctx context.Context) http.Handler {
	return admin.useraction(ctx, "toggle admin", func(r *http.Request, user *auth.AuthUser) error {
		if user.IsAdmin() {
			return user.MakeRegular(r.Context())
		}
		return user.MakeAdmin(r.Context())
	})
}

func (admin *Admin) usersPasswordHandler(ctx context.Context) http.Handler {
	return admin.useraction(ctx, "set password", func(r *http.Request, user *auth.AuthUser) error {
		password := r.PostFormValue("password")
		if password == "" {
			return httpx.ErrBadRequest
		}
		// check the password policy
		err := user.CheckPasswordPolicy(password)
		if err != nil {
			return err
		}
		return user.SetPassword(r.Context(), []byte(password))
	})
}
