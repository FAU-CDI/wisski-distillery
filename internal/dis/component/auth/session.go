package auth

import (
	"context"
	"html/template"
	"net/http"
	"net/url"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
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
	name, ok := sess.Values[sessionUserKey]
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
	if !user.Enabled {
		return nil, nil
	}

	// get the user
	return user, nil
}

const sessionCookieName = "distillery-session"

// session returns the session that belongs to a given request.
// If the session is not set, creates a new session.
func (auth *Auth) session(r *http.Request) (*sessions.Session, error) {
	return auth.store.Get(func() sessions.Store {
		return sessions.NewCookieStore([]byte(auth.Config.SessionSecret))
	}).Get(r, sessionCookieName)
}

const sessionUserKey = "user"

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
	sess.Values[sessionUserKey] = user.User.User
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

//go:embed "templates/login.html"
var loginHTMLStr string
var loginTemplate = static.AssetsAuthLogin.MustParseShared("login.html", loginHTMLStr)

var loginResponse = httpx.Response{
	ContentType: "text/plain",
	Body:        []byte("user is signed in"),
}

type authloginContext struct {
	Message string
	Form    template.HTML
}

// Protect returns a new handler which requires a user to be logged in and pass the perm function.
//
// If an unauthenticated user attempts to access the returned handler, they are redirected to the login endpoint.
// When a user is logged in, and they pass the perm function (or the perm function is nil), the original handler is called.
func (auth *Auth) Protect(handler http.Handler, perm func(user *AuthUser, r *http.Request) (ok bool, err error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// load the user in the session
		user, err := auth.UserOf(r)
		if err != nil {
			goto err
		}

		// if there is no user in the session
		// we need to login the user
		if user == nil {
			// we can't redirect anything other than GET
			// (because it might be a form)
			// => so we just return a forbidden
			if r.Method != http.MethodGet {
				goto forbidden
			}

			// redirect the user to the login endpoint, with the original URI as a return
			dest := "/auth/login?next=" + url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, dest, http.StatusSeeOther)
			return
		}

		// if we have a permission check, we need to call it
		// to find out if the user is actually allowed to access the page
		if perm != nil {
			ok, err := perm(user, r)
			if err != nil {
				goto err
			}
			if !ok {
				goto forbidden
			}
		}

		// store the user into the session
		r = r.WithContext(context.WithValue(r.Context(), ctxUserKey, user))
		handler.ServeHTTP(w, r)
		return
	forbidden:
		httpx.HTMLInterceptor.Intercept(w, r, httpx.ErrForbidden)
		return
	err:
		httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
	})
}

// authLogin implements a view to login a user
func (auth *Auth) authLogin(ctx context.Context) http.Handler {
	return &httpx.Form[*AuthUser]{
		Fields: []httpx.Field{
			{Name: "username", Type: httpx.TextField, Label: "Username"},
			{Name: "password", Type: httpx.PasswordField, EmptyOnError: true, Label: "Password"},
			{Name: "passcode", Type: httpx.TextField, EmptyOnError: true, Label: "Passcode (optional)"},
		},
		FieldTemplate: httpx.PureCSSFieldTemplate,

		CSRF: auth.csrf.Get(nil),

		RenderForm: func(template template.HTML, err error, w http.ResponseWriter, r *http.Request) {
			ctx := authloginContext{
				Message: "",
				Form:    template,
			}
			if err != nil {
				ctx.Message = "Login Failed"
			}
			httpx.WriteHTML(ctx, nil, loginTemplate, "", w, r)
		},

		Validate: func(r *http.Request, values map[string]string) (*AuthUser, error) {
			username, password, passcode := values["username"], values["password"], values["passcode"]

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
