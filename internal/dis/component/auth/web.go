package auth

import (
	"context"
	"net/http"
	"net/url"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"

	_ "embed"
)

func (auth *Auth) Routes() []string {
	return []string{"/auth/"}
}

type contextUserKey struct{}

var ctxUserKey = contextUserKey{}

const (
	sessionCookieName = "distillery-session"
	sessionUserKey    = "user"
)

// session returns the session belonging to a request
func (auth *Auth) session(r *http.Request) (*sessions.Session, error) {
	auth.storeOnce.Do(func() {
		auth.store = sessions.NewCookieStore([]byte(auth.Config.SessionSecret))
	})
	return auth.store.Get(r, sessionCookieName)
}

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

// writeLogin marks the user as logged in on the given writer
func (auth *Auth) writeLogin(w http.ResponseWriter, r *http.Request, user *AuthUser) error {
	sess, err := auth.session(r)
	if err != nil {
		return err
	}
	sess.Values[sessionUserKey] = user.User.User
	return sess.Save(r, w)
}

// writeLogout logs out the user form the given session
func (auth *Auth) writeLogout(w http.ResponseWriter, r *http.Request) error {
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

// HandleRoute returns the handler for the requested route
func (auth *Auth) HandleRoute(ctx context.Context, route string) (http.Handler, error) {
	router := httprouter.New()

	router.Handler(http.MethodGet, route, auth.Protect(loginResponse, nil))

	router.HandlerFunc(http.MethodGet, route+"login", auth.loginRoute)
	router.HandlerFunc(http.MethodPost, route+"login", auth.loginRoute)

	router.HandlerFunc(http.MethodGet, route+"logout", auth.logoutRoute)

	return router, nil
}

type loginContext struct {
	Message string
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

func (auth *Auth) loginRoute(w http.ResponseWriter, r *http.Request) {
	var message string

	// try to read a user from the session
	user, err := auth.UserOf(r)
	if err != nil {
		httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
		return
	}

	if user != nil {
		goto success
	}

	switch r.Method {
	default:
		panic("never reached")
	case http.MethodGet:
		goto form
	case http.MethodPost:
		// parse the form!
		if err := r.ParseForm(); err != nil {
			message = "Login failed"
			goto form
		}

		// get the username and password
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		// make sure that the user exists
		user, err := auth.User(r.Context(), username)
		if err != nil {
			message = "Login failed"
			goto form
		}

		// check the password (TODO: Support TOTP)
		err = user.CheckPassword(r.Context(), []byte(password))
		if err != nil {
			message = "Login failed"
			goto form
		}

		// and we logged the user in!
		auth.writeLogin(w, r, user)
		goto success
	}

form:
	httpx.WriteHTML(loginContext{
		Message: message,
	}, nil, loginTemplate, "", w, r)
	return
success:
	// get the destination
	next := r.URL.Query().Get("next")
	if next == "" || next[0] != '/' {
		next = "/"
	}

	// and redirect to it!
	http.Redirect(w, r, next, http.StatusSeeOther)
}

func (auth *Auth) logoutRoute(w http.ResponseWriter, r *http.Request) {
	// do the logout
	auth.writeLogout(w, r)

	// get the destination
	next := r.URL.Query().Get("next")
	if next == "" || next[0] != '/' {
		next = "/"
	}

	// and redirect to it!
	http.Redirect(w, r, next, http.StatusSeeOther)

}
