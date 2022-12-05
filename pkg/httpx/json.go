package httpx

import (
	"encoding/json"
	"net/http"
)

// JSON creates a new JSONHandler
func JSON[T any](f func(r *http.Request) (T, error)) JSONHandler[T] {
	return JSONHandler[T](f)
}

// WriteJSON writes a JSON response of type T to w.
// If an error occured, writes an error response instead.
func WriteJSON[T any](result T, err error, w http.ResponseWriter, r *http.Request) {
	// handle any errors
	if JSONInterceptor.Intercept(w, r, err) {
		return
	}

	// write out the response as json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// JSONHandler implements [http.Handler] by returning values as json to the caller.
// In case of an error, a generic "internal server error" message is returned.
type JSONHandler[T any] func(r *http.Request) (T, error)

// ServeHTTP calls j(r) and returns json
func (j JSONHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result, err := j(r)
	WriteJSON(result, err, w, r)
}
