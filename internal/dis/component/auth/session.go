package auth

import (
	"context"
	"errors"
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/tkw1536/pkglib/httpx"
	"github.com/tkw1536/pkglib/httpx/field"

	"github.com/gorilla/sessions"

	_ "embed"
)

// SessionOf returns the session and user logged into the provided request.
// token indicates if the user used a token to authenticate, or a browser session was used.
// A token takes priority over a user in a session.
//
// If there is no user associated with the given request, user and error are nil, and token is false.
// An invalid session, expired token, or disabled user all result in user = nil.
//
// When no SessionOf exists in the given session returns nil.
func (auth *Auth) SessionOf(r *http.Request) (session component.SessionInfo, user *AuthUser, err error) {
	// check the user from the token first
	{
		user, err := auth.UserOfToken(r)
		if user != nil && err == nil {
			return component.SessionInfo{User: &user.User, Token: true}, user, nil
		}
	}

	// fallback to using session
	{
		user, err := auth.UserOfSession(r)
		if err != nil {
			return component.SessionInfo{}, nil, err
		}
		if user == nil {
			return component.SessionInfo{}, nil, nil
		}
		return component.SessionInfo{User: &user.User, Token: false}, user, nil
	}
}

// UserOfToken returns the user associated with the token in request.
// To check the user of a token or session, use SessionOf.
func (auth *Auth) UserOfToken(r *http.Request) (user *AuthUser, err error) {
	// get the token object
	token, err := auth.dependencies.Tokens.TokenOf(r)
	if token == nil {
		return nil, err
	}
	return auth.checkUser(r.Context(), token.User)
}

// UserOfSession returns the user of the session associated with r.
func (auth *Auth) UserOfSession(r *http.Request) (user *AuthUser, err error) {
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
	name, ok := sess.Values[server.SessionUserKey]
	if !ok {
		return nil, nil
	}
	nameS, ok := name.(string)
	if !ok || nameS == "" {
		return nil, nil
	}
	return auth.checkUser(ctx, nameS)
}

func (auth *Auth) checkUser(ctx context.Context, name string) (user *AuthUser, err error) {
	// fetch the user, check if they still exist
	user, err = auth.User(ctx, name)
	if err == ErrUserNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// user isn't enabled
	if user == nil || !user.IsEnabled() {
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
	}).Get(r, server.SessionCookie)
}

func (auth *Auth) Menu(r *http.Request) []component.MenuItem {

	user, err := auth.UserOfSession(r)
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
	sess.Values[server.SessionUserKey] = user.User.User
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
var loginTemplate = templating.ParseForm(
	"login.html", loginHTML, httpx.FormTemplate,

	templating.Title("Login Required"),
	templating.Assets(assets.AssetsUser),
)

var errLoginFailed = errors.New("login failed")

// authLogin implements a view to login a user
func (auth *Auth) authLogin(ctx context.Context) http.Handler {
	tpl := loginTemplate.Prepare(
		auth.dependencies.Templating,
		func(flags templating.Flags, r *http.Request) templating.Flags {
			flags.Crumbs = []component.MenuItem{
				{Title: "Login", Path: template.URL(r.URL.RequestURI())},
			}
			return flags
		},
	)

	return &httpx.Form[*AuthUser]{
		Fields: []field.Field{
			{Name: "username", Type: field.Text, Autocomplete: field.Username, Label: "Username"},
			{Name: "password", Type: field.Password, Autocomplete: field.CurrentPassword, EmptyOnError: true, Label: "Password"},
			{Name: "otp", Type: field.Text, Autocomplete: field.OneTimeCode, EmptyOnError: true, Label: "Passcode (optional)"},
		},
		FieldTemplate: field.PureCSSFieldTemplate,

		RenderTemplateContext: func(ctx httpx.FormContext, r *http.Request) any {
			if ctx.Err != nil {
				ctx.Err = errLoginFailed
			}
			return tpl.Context(r, templating.NewFormContext(ctx))
		},
		RenderTemplate: tpl.Template(),

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
			user, err := auth.UserOfSession(r)
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
