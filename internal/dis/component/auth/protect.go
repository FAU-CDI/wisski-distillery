package auth

import (
	"context"
	"net/http"
	"net/url"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/tkw1536/pkglib/httpx"
)

// Protect returns a new handler which requires a user to be logged in and have the provided scope.
//
// AllowToken determines if a token is allowed instead of a user session.
//
// If an unauthenticated user attempts to access the returned handler, they are redirected to the login endpoint.
// If an authenticated user is missing the given scope, a Forbidden response is called.
// If an authenticated calls the endpoint, and they have the given permissions, the original handler is called.
func (auth *Auth) Protect(handler http.Handler, AllowToken bool, scope component.Scope, param func(*http.Request) string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var forbiddenMessage string
		var paramValue string

		// load the user in the session
		// TODO: In a future version of sessions, check if token has the permitted scope.
		session, user, err := auth.SessionOf(r)
		if err != nil {
			goto err
		}

		// token was set, but not allowed!
		if session.Token && !AllowToken {
			goto forbidden
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

		// check if we need to load the parameter
		if param != nil {
			paramValue = param(r)
		}

		// check if we have the actual scope
		{
			err = auth.CheckScope(paramValue, scope, r)
			if ade, ok := err.(component.AccessDeniedError); ok {
				forbiddenMessage = ade.Error()
				goto forbidden
			}
			if err != nil {
				goto err
			}
		}

		// store the user into the session, and then return the new session
		r = r.WithContext(context.WithValue(r.Context(), ctxUserKey, user))
		handler.ServeHTTP(w, r)
		return
	forbidden:
		{
			httpx.Response{
				ContentType: "text/plain",
				StatusCode:  http.StatusForbidden,
				Body:        []byte(forbiddenMessage),
			}.ServeHTTP(w, r)
			return
		}
	err:
		httpx.TextInterceptor.Fallback.ServeHTTP(w, r)
	})
}

// Require returns a slice containing one decorator that acts like auth.Protect(allowToken,scope,param) on every request.
func (auth *Auth) Require(allowToken bool, scope component.Scope, param func(*http.Request) string) func(http.Handler) http.Handler {
	// TODO: Work on this stuff
	return func(h http.Handler) http.Handler {
		return auth.Protect(h, allowToken, scope, param)
	}
}
