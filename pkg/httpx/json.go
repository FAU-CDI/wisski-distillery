package httpx

import (
	"encoding/json"
	"net/http"
)

// JSON creates a new JSONHandler
func JSON[T any](f func(r *http.Request) (T, error)) JSONHandler[T] {
	return JSONHandler[T](f)
}

// JSONHandler implements [http.Handler] by returning values as json to the caller.
// In case of an error, a generic "internal server error" message is returned.
type JSONHandler[T any] func(r *http.Request) (T, error)

// ServeHTTP calls j(r) and returns json
func (j JSONHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// call the function
	result, err := j(r)

	// handle any errors
	if jsonInterceptor.Intercept(w, r, err) {
		return
	}

	// write out the response as json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
