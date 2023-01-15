package httpx

import (
	"context"
	"net/http"
)

// WithContextWrapper generates a new handler that wraps the context of each request with the wrapper function.
func WithContextWrapper(handler http.Handler, wrapper func(context.Context) context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(wrapper(r.Context()))
		handler.ServeHTTP(w, r)
	})
}
