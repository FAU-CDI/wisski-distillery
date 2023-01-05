package component

import (
	"context"
	"net/http"
)

// Routeable is a component that is servable
type Routeable interface {
	Component

	// Routes returns information about the routes to be handled by this Routeable
	Routes() Routes

	// HandleRoute returns the handler for the requested path
	HandleRoute(ctx context.Context, path string) (http.Handler, error)
}

// Routes represents information about a single Routeable
type Routes struct {
	// Paths are the paths handled by this routeable.
	// Each path is passed to HandleRoute() individually.
	Paths []string

	// CSRF indicates if this route should be protected by CSRF.
	// CSRF protection is applied prior to any custom decorator being called.
	CSRF bool

	// Decorators is a function applied to the handler returned by HandleRoute.
	// When nil, it is not applied.
	Decorator func(http.Handler) http.Handler
}

// Decorate decorates the provided handler with the options specified in this handler.
func (routes Routes) Decorate(handler http.Handler, csrf func(http.Handler) http.Handler) http.Handler {
	if routes.CSRF && csrf != nil {
		handler = csrf(handler)
	}
	if routes.Decorator != nil {
		handler = routes.Decorator(handler)
	}
	return handler
}
