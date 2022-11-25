package httpx

import (
	"net/http"
	"text/template"
)

// RedirectHandler represents a handler that redirects the user to the address returned
type RedirectHandler func(r *http.Request) (string, int, error)

// ServeHTTP calls r(r) and returns json
func (rh RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// call the function
	url, code, err := rh(r)

	// intercept the errors
	if textInterceptor.Intercept(w, r, err) {
		return
	}

	// do the redirect
	http.Redirect(w, r, url, code)
}

type ClientSideRedirect func(r *http.Request) (string, error)

var htmlTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html lang="en">
<title>Redirecting</title>
<meta http-equiv="refresh" content="0; url={{ . }}" />
You should be redirected to <a href="{{ . }}">{{ . }}</a>
`))

// ServeHTTP calls r(r) and returns json
func (rh ClientSideRedirect) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// call the function
	url, err := rh(r)

	// intercept the errors
	if htmlInterceptor.Intercept(w, r, err) {
		return
	}

	// write out the response as json
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	htmlTemplate.Execute(w, url)
}
