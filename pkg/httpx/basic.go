package httpx

import "net/http"

var basicUnauthorized = []byte("Unauthorized")

// BasicAuth returns a new [http.Handler] that requires any credentials to pass the check function
func BasicAuth(handler http.Handler, realm string, check func(username, password string) bool) http.Handler {
	var authenticateHeader = `Basic realm="` + realm + `"`
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if the basic authentication passes
		// we can just use the handler!
		user, pass, ok := r.BasicAuth()
		if ok && check(user, pass) {
			handler.ServeHTTP(w, r)
			return
		}

		// http authentication did not pass
		w.Header().Add("WWW-Authenticate", authenticateHeader)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(basicUnauthorized)
	})
}
