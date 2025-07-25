//spellchecker:words admin
package admin

//spellchecker:words context errors http github wisski distillery internal component auth server assets templating wdlog pkglib httpx form field embed
import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"go.tkw01536.de/pkglib/httpx"
	"go.tkw01536.de/pkglib/httpx/form"
	"go.tkw01536.de/pkglib/httpx/form/field"

	_ "embed"
)

//go:embed "html/users.html"
var usersHTML []byte
var usersTemplate = templating.Parse[usersContext](
	"users.html", usersHTML, nil,

	templating.Title("Users"),
	templating.Assets(assets.AssetsAdmin),
)

type usersContext struct {
	templating.RuntimeFlags
	Error string
	Users []*auth.AuthUser
}

func (admin *Admin) users(context.Context) http.Handler {
	tpl := usersTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuUsers,
		),
		templating.Actions(
			menuUserCreate,
		),
	)

	return tpl.HTMLHandler(admin.dependencies.Handling, func(r *http.Request) (uc usersContext, err error) {
		uc.Error = r.URL.Query().Get("error")
		uc.Users, err = admin.dependencies.Auth.Users(r.Context())
		return
	})
}

//go:embed "html/user_create.html"
var userCreateHTML []byte
var userCreateTemplate = templating.ParseForm(
	"user_create.html", userCreateHTML, form.FormTemplate,

	templating.Title("Create User"),
	templating.Assets(assets.AssetsAdmin),
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

func (admin *Admin) createUser(context.Context) http.Handler {
	tpl := userCreateTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuUsers,
			menuUserCreate,
		),
	)

	return &form.Form[createUserResult]{
		Fields: []field.Field{
			{Name: "username", Type: field.Text, Autocomplete: field.Username, Label: "Username"},
			{Name: "password", Type: field.Password, Autocomplete: field.NewPassword, Label: "Password"},
			{Name: "admin", Type: field.Checkbox, Label: "Distillery Administrator"},
		},
		FieldTemplate: assets.PureCSSFieldTemplate,

		Template:         tpl.Template(),
		TemplateContext:  templating.FormTemplateContext(tpl),
		LogTemplateError: tpl.LogTemplateError,

		Validate: func(r *http.Request, values map[string]string) (cu createUserResult, err error) {
			cu.User, cu.Passsword, cu.Admin = values["username"], values["password"], values["admin"] == field.CheckboxChecked

			if cu.User == "" {
				return cu, errCreateInvalidUsername
			}
			if cu.Passsword == "" {
				return cu, errCreateInvalidPassword
			}

			// check the password policy
			err = admin.dependencies.Auth.CheckPasswordPolicy(cu.Passsword, cu.User)
			if err != nil {
				return cu, fmt.Errorf("failed to check passsword policy: %w", err)
			}

			return cu, nil
		},

		Success: func(cu createUserResult, values map[string]string, w http.ResponseWriter, r *http.Request) error {
			// create the user
			user, err := admin.dependencies.Auth.CreateUser(r.Context(), cu.User)
			if err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}

			// disable the user and setup the admin flag
			user.SetAdmin(cu.Admin)
			if err := user.Save(r.Context()); err != nil {
				return fmt.Errorf("failed to save user: %w", err)
			}

			// set the password!
			err = user.SetPassword(r.Context(), []byte(cu.Passsword))
			if err != nil {
				return fmt.Errorf("failed to set password: %w", err)
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
	logger := wdlog.Of(ctx)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			logger.Error(
				"failed to parse form",
				"error", err,
				"action", name,
			)
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		username := r.PostFormValue("user")
		user, err := admin.dependencies.Auth.User(r.Context(), username)
		if err != nil {
			logger.Error(
				"failed to get user",
				"error", err,
				"action", name,
			)
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		me, err := admin.dependencies.Auth.UserOfSession(r)
		if err != nil {
			logger.Error(
				"failed to get current user",
				"error", err,
				"action", name,
			)
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		// don't allow the current user
		if me.User.User == user.User.User {
			errNotCurrentUser.ServeHTTP(w, r)
			return
		}

		if err := action(r, user); err != nil {
			logger.Error(
				"failed to act on user",
				"error", err,
				"action", name,
			)
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
			return fmt.Errorf("failed to check password policy: %w", err)
		}
		return user.SetPassword(r.Context(), []byte(password))
	})
}

func (admin *Admin) usersUnsetPasswordHandler(ctx context.Context) http.Handler {
	return admin.useraction(ctx, "unset password", func(r *http.Request, user *auth.AuthUser) error {
		user.PasswordHash = nil
		return user.Save(r.Context())
	})
}

func (admin *Admin) usersImpersonateHandler(ctx context.Context) http.Handler {
	logger := wdlog.Of(ctx)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			logger.Error(
				"failed to parse form",
				"error", err,
				"action", "impersonate",
			)
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		username := r.PostFormValue("user")
		user, err := admin.dependencies.Auth.User(r.Context(), username)
		if err != nil {
			logger.Error(
				"failed to get user",
				"error", err,
				"action", "impersonate",
			)

			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		// login the user into the session of the provided user
		if err := admin.dependencies.Auth.Login(w, r, user); err != nil {
			logger.Error(
				"failed to login user",
				"error", err,
				"action", "impersonate",
			)
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		// and go there
		http.Redirect(w, r, "/user/", http.StatusSeeOther)
	})
}
