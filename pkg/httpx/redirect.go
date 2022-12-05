package httpx

import (
	"net/http"
)

// RedirectHandler represents a handler that redirects the user to the address returned
type RedirectHandler func(r *http.Request) (string, int, error)

// ServeHTTP calls r(r) and returns json
func (rh RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// call the function
	url, code, err := rh(r)

	// intercept the errors
	if TextInterceptor.Intercept(w, r, err) {
		return
	}

	// do the redirect
	http.Redirect(w, r, url, code)
}
