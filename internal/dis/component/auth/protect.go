package auth

import (
	"context"
	"net/http"
	"net/url"

	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
)

// Protect returns a new handler which requires a user to be logged in and pass the perm function.
//
// If an unauthenticated user attempts to access the returned handler, they are redirected to the login endpoint.
// If an authenticated user is missing permissions, a Forbidden response is called.
// If an authenticated calls the endpoint, and they have the given permissions, the original handler is called.
func (auth *Auth) Protect(handler http.Handler, perm Permission) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var grant Grant

		// load the user in the session
		user, err := auth.UserOf(r)
		if err != nil {
			goto err
		}

		// if there is no user in the session, they need to login first!
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

		{
			var err error
			// call the permission check
			grant, err = perm.Permit(user, r)
			if err != nil {
				goto err
			}
			if !grant.Granted() {
				goto forbidden
			}
		}

		// store the user into the session, and then return the new session
		r = r.WithContext(context.WithValue(r.Context(), ctxUserKey, user))
		handler.ServeHTTP(w, r)
		return
	forbidden:
		{
			message := "Forbidden"
			if grant != nil {
				message = grant.Denied()
			}
			httpx.Response{
				ContentType: "text/plain",
				StatusCode:  http.StatusForbidden,
				Body:        []byte(message),
			}.ServeHTTP(w, r)
			return
		}
	err:
		httpx.TextInterceptor.Fallback.ServeHTTP(w, r)
	})
}

// Require returns a slice containing one decorator that acts like Protect(perm) on every request.
// It returns
func (auth *Auth) Require(perm Permission) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return auth.Protect(h, perm)
	}
}

// Has checks if the given request has the given permission.
// If an error occurs, returns false.
func (auth *Auth) Has(perm Permission, r *http.Request) bool {
	user, err := auth.UserOf(r)
	if err != nil || user == nil {
		return false
	}
	ok, err := perm.Permit(user, r)
	return err == nil && ok.Granted()
}

// Admin represents a permission that checks if a user is an administrator and has totp enabled.
var Admin Permission = func(user *AuthUser, r *http.Request) (ok Grant, err error) {
	return Bool2Grant(user != nil && user.IsAdmin() && user.IsTOTPEnabled(), "user needs to have admin permissions and passcode enabled"), nil
}

// User represents a permission that checks if a user is enabled
var User Permission = func(user *AuthUser, r *http.Request) (ok Grant, err error) {
	return Bool2Grant(user != nil && user.IsEnabled(), "user needs to be enabled"), nil
}
