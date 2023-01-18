package auth

import (
	"context"
	"errors"
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static/custom"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx/field"
	"github.com/gorilla/sessions"

	_ "embed"
)

// UserOf returns the user logged into the given request.
// If there is no user associated with the given user, user and error will be nil.
//
// When no UserOf exists in the given session returns nil.
// An invalid session (for a UserOf)
func (auth *Auth) UserOf(r *http.Request) (user *AuthUser, err error) {
	ctx := r.Context()
	if user, ok := ctx.Value(ctxUserKey).(*AuthUser); ok && user != nil {
		return user, nil
	}

	// first read the session
	sess, err := auth.session(r)
	if err != nil {
		return nil, err
	}

	// try to read the name from the session
	name, ok := sess.Values[control.SessionUserKey]
	if !ok {
		return nil, nil
	}
	nameS, ok := name.(string)
	if !ok || nameS == "" {
		return nil, nil
	}

	// fetch the user, check if they still exist
	user, err = auth.User(ctx, nameS)
	if err == ErrUserNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// user isn't enabled
	if !user.IsEnabled() {
		return nil, nil
	}

	// get the user
	return user, nil
}

// session returns the session that belongs to a given request.
// If the session is not set, creates a new session.
func (auth *Auth) session(r *http.Request) (*sessions.Session, error) {
	return auth.store.Get(func() sessions.Store {
		return sessions.NewCookieStore([]byte(auth.Config.SessionSecret))
	}).Get(r, control.SessionCookie)
}

func (auth *Auth) Menu(r *http.Request) []component.MenuItem {

	user, err := auth.UserOf(r)
	if user == nil || err != nil {
		return nil
	}
	return []component.MenuItem{
		{
			Title:    "Logout",
			Path:     "/auth/logout",
			Priority: component.MenuAuth,
		},
	}

}

type contextUserKey struct{}

var ctxUserKey = contextUserKey{}

// Login logs a user into the given request.
//
// If a user was previously logged into this session,
// UserOf may not return the correct user until the user makes a new request.
//
// It is recommended to send a HTTP redirect to make sure a new request is made.
func (auth *Auth) Login(w http.ResponseWriter, r *http.Request, user *AuthUser) error {
	sess, err := auth.session(r)
	if err != nil {
		return err
	}
	sess.Values[control.SessionUserKey] = user.User.User
	return sess.Save(r, w)
}

// Logout logs out the user from the given session.
//
// UserOf may return incorrect results until the user makes a new request.
// It is recommended to send a HTTP redirect to make sure a new request is made.
func (auth *Auth) Logout(w http.ResponseWriter, r *http.Request) error {
	sess, err := auth.session(r)
	if err != nil {
		return err
	}
	sess.Options.MaxAge = -1
	return sess.Save(r, w)
}

//go:embed "login.html"
var loginHTML []byte
var loginTemplate = custom.ParseForm("login.html", loginHTML, static.AssetsUser)

var loginResponse = httpx.Response{
	ContentType: "text/plain",
	Body:        []byte("user is signed in"),
}

var errLoginFailed = errors.New("Login failed")

// authLogin implements a view to login a user
func (auth *Auth) authLogin(ctx context.Context) http.Handler {
	tpl := loginTemplate.Prepare(auth.Dependencies.Custom)

	return &httpx.Form[*AuthUser]{
		Fields: []field.Field{
			{Name: "username", Type: field.Text, Autocomplete: field.Username, Label: "Username"},
			{Name: "password", Type: field.Password, Autocomplete: field.CurrentPassword, EmptyOnError: true, Label: "Password"},
			{Name: "otp", Type: field.Text, Autocomplete: field.OneTimeCode, EmptyOnError: true, Label: "Passcode (optional)"},
		},
		FieldTemplate: field.PureCSSFieldTemplate,

		RenderForm: func(context httpx.FormContext, w http.ResponseWriter, r *http.Request) {
			if context.Err != nil {
				context.Err = errLoginFailed
			}
			tpl.Execute(w, r, custom.BaseFormContext{FormContext: context}, custom.BaseContextGaps{
				Crumbs: []component.MenuItem{
					{Title: "Login", Path: template.URL(r.URL.RequestURI())},
				},
			})
		},

		Validate: func(r *http.Request, values map[string]string) (*AuthUser, error) {
			username, password, passcode := values["username"], values["password"], values["otp"]

			// make sure that the user exists
			user, err := auth.User(ctx, username)
			if err != nil {
				return nil, err
			}

			// check the password and totp
			err = user.CheckCredentials(ctx, []byte(password), passcode)
			if err != nil {
				return nil, err
			}
			return user, nil
		},

		SkipForm: func(r *http.Request) (user *AuthUser, skip bool) {
			user, err := auth.UserOf(r)
			return user, err == nil && user != nil
		},

		RenderSuccess: func(user *AuthUser, _ map[string]string, w http.ResponseWriter, r *http.Request) error {
			if err := auth.Login(w, r, user); err != nil {
				return err
			}

			// get the destination
			next := r.URL.Query().Get("next")
			if next == "" || next[0] != '/' {
				next = "/"
			}

			// and redirect to it!
			http.Redirect(w, r, next, http.StatusSeeOther)

			return nil
		},
	}
}

// authLogout implements the authLogout view to logout a user
func (auth *Auth) authLogout(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do the logout
		auth.Logout(w, r)

		// get the destination
		next := r.URL.Query().Get("next")
		if next == "" || next[0] != '/' {
			next = "/"
		}

		// and redirect to it!
		http.Redirect(w, r, next, http.StatusSeeOther)
	})
}
