package auth

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
)

// Permission represents a permission granted to a user.
//
// The nil permission represents any authenticated user.
type Permission func(user *AuthUser, r *http.Request) (ok Grant, err error)

// Grant represents an object that either grants or denies access for a certain permission
type Grant interface {
	isGranted()

	// Granted returns a boolean indicating if permission to the resource in question
	// has been granted
	Granted() bool

	// Denied returns a string containing an error message to display to the user when permission is denied.
	// When Granted() returns true, the behaviour is undefined.
	Denied() string
}

// Bool2Grant returns a new grant that returns granted for the given boolean, and message as the denied message.
func Bool2Grant(granted bool, message string) Grant {
	if granted {
		return grantAllow{}
	}
	return grantDeny(message)
}

type grantAllow struct{}

func (grantAllow) isGranted()     {}
func (grantAllow) Granted() bool  { return true }
func (grantAllow) Denied() string { return "" }

type grantDeny string

func (grantDeny) isGranted()      {}
func (g grantDeny) Granted() bool { return false }
func (g grantDeny) Denied() string {
	if g == "" {
		return "Forbidden"
	}
	return string(g)
}

var errPermissionPanic = errors.New("permission: panic()")

// Permit checks if the given user has this permission.
func (perm Permission) Permit(user *AuthUser, r *http.Request) (ok Grant, err error) {
	// if there is no permission, then we just check if there is some user
	if perm == nil {
		return Bool2Grant(user != nil, ""), nil
	}

	// recover any panic()ed permission call
	// to prevent the handler from panic()ing
	defer func() {
		if p := recover(); p != nil {
			ok = Bool2Grant(false, "unknown error")
			err = errPermissionPanic
		}
	}()

	return perm(user, r)
}

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

// Admin represents a permission that checks if a user is an administrator and has totp enabled.
var Admin Permission = func(user *AuthUser, r *http.Request) (ok Grant, err error) {
	return Bool2Grant(user != nil && user.Admin && user.TOTPEnabled, "user needs to have admin permissions and TOTP enabled"), nil
}
