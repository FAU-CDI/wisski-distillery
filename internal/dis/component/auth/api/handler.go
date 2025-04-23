// Package api implements a common handler used by the api routes
package api

//spellchecker:words encoding json errors http github wisski distillery internal config component auth scopes wdlog pkglib httpx lazy
import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/tkw1536/pkglib/httpx"
	"github.com/tkw1536/pkglib/lazy"
)

// Handler represents an API handler that returns a REST response.
// The response is automatically marshaled using T.
type Handler[T any] struct {
	Config *config.Config
	Auth   *auth.Auth // Handler to handle Auth

	Methods []string // HTTP methods to allow
	methods lazy.Lazy[map[string]struct{}]

	Scope      scopes.Scope
	ScopeParam func(*http.Request) string
	Handler    func(string, *http.Request) (T, error)
}

var apiNotEnabled = &Response{
	Status:  http.StatusNotImplemented,
	Message: "API is not implemented on this server",
}

var apiMethodNotAllowed = &Response{
	Status:  http.StatusMethodNotAllowed,
	Message: "method not allowed",
}

var apiInternalServerError = &Response{
	Status:  http.StatusInternalServerError,
	Message: "internal server error",
}

var apiBadRequest = &Response{
	Status:  http.StatusBadRequest,
	Message: "bad request",
}

var apiNotFound = &Response{
	Status:  http.StatusNotFound,
	Message: "not found",
}

var apiForbidden = &Response{
	Status:  http.StatusForbidden,
	Message: "forbidden",
}

// ServeHTTP servers an api call.
func (handler *Handler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check that the api is actually enabled!
	if !handler.Config.HTTP.API.Value {
		apiNotEnabled.ServeHTTP(w, r)
		return
	}

	// get the permitted methods
	methods := handler.methods.Get(func() map[string]struct{} {
		m := make(map[string]struct{}, len(handler.Methods)+1)
		for _, method := range handler.Methods {
			m[method] = struct{}{}
		}
		m["OPTIONS"] = struct{}{}
		return m
	})

	// check that the method is permitted
	if _, ok := methods[r.Method]; !ok {
		apiMethodNotAllowed.ServeHTTP(w, r)
		return
	}

	// we now delegate to user-level code;
	// so we now need to make sure that panic()s are caught.
	var stage string

	//nolint:contextcheck
	defer func() {
		// recover any error
		rec := recover()
		if rec == nil {
			return
		}

		// log the error, and serve the default internal server error
		wdlog.Of(r.Context()).Error(
			"api handler caused panic()",
			"panic", fmt.Sprint(rec),
			"stage", stage,
			"route", r.URL.RequestURI(),
		)
		apiInternalServerError.ServeHTTP(w, r)
	}()

	// read the parameter
	stage = "param"
	var param string
	if handler.ScopeParam != nil {
		param = handler.ScopeParam(r)
	}

	// check that the scope is correct
	stage = "scope"
	if err := handler.Auth.CheckScope(param, handler.Scope, r); err != nil {
		(&Response{
			Status:  http.StatusForbidden,
			Message: err.Error(),
		}).ServeHTTP(w, r)
		return
	}

	stage = "handler"

	result, err := handler.Handler(param, r)
	switch {
	case err == nil: /* keep going */

	// handle common httpx errors
	case errors.Is(err, httpx.ErrInternalServerError):
		apiInternalServerError.ServeHTTP(w, r)
		return
	case errors.Is(err, httpx.ErrBadRequest):
		apiBadRequest.ServeHTTP(w, r)
		return
	case errors.Is(err, httpx.ErrNotFound):
		apiNotFound.ServeHTTP(w, r)
		return
	case errors.Is(err, httpx.ErrForbidden):
		apiForbidden.ServeHTTP(w, r)
		return
	case errors.Is(err, httpx.ErrMethodNotAllowed):
		apiMethodNotAllowed.ServeHTTP(w, r)
		return

		// generic error
	default:
		(&Response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}).ServeHTTP(w, r)
		return
	}

	stage = "marshal"

	// encode the result into json and send it as the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// Response objects cache response serialization.
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	res     lazy.Lazy[httpx.Response]
}

func (g *Response) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.res.Get(func() httpx.Response {
		bytes, _ := json.Marshal(g)
		return httpx.Response{
			ContentType: "application/json",
			Body:        bytes,
			StatusCode:  g.Status,
		}
	}).ServeHTTP(w, r)
}
