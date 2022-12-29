package auth

import (
	"context"
	"net/http"
	"net/url"

	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
)

// Permission represents a permission for a user
//
// The nil permission represents any authenticated user.
type Permission func(user *AuthUser, r *http.Request) (ok bool, err error)

// Protect returns a new handler which requires a user to be logged in and pass the perm function.
//
// If an unauthenticated user attempts to access the returned handler, they are redirected to the login endpoint.
// If an authenticated user is missing permissions, a Forbidden response is called.
// If an authenticated calls the endpoint, and they have the given permissions, the original handler is called.
func (auth *Auth) Protect(handler http.Handler, perm Permission) http.Handler {
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

// Admin represents a permission that checks if a user is an administrator and has totp enabled.
var Admin Permission = func(user *AuthUser, r *http.Request) (ok bool, err error) {
	return user != nil && user.Admin && user.TOTPEnabled, nil
}
