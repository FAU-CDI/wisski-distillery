package component

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/pkg/mux"
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
	// Prefix is the prefix this pattern handles
	Prefix string

	// MatchAllDomains indicates that all domains, even the non-default domain, should be matched
	MatchAllDomains bool

	// MenuTitle and MenuPriority return the priority and title of this menu item
	MenuTitle    string
	MenuPriority MenuPriority

	// Exact indicates that only the exact prefix, as opposed to any sub-paths, are matched.
	// Trailing '/'s are automatically trimmed, even with an exact match.
	Exact bool

	// Aliases are the additional prefixes this route handles.
	Aliases []string

	// CSRF indicates if this route should be protected by CSRF.
	// CSRF protection is applied prior to any custom decorator being called.
	CSRF bool

	// Decorators is a function applied to the handler returned by HandleRoute.
	// When nil, it is not applied.
	Decorator func(http.Handler) http.Handler
}

type RouteContext struct {
	DefaultDomain bool
}

// Predicate returns the predicate corresponding to the given route
func (routes Routes) Predicate(context func(*http.Request) RouteContext) mux.Predicate {
	if routes.MatchAllDomains {
		return nil
	}

	// match only the default domain
	return func(r *http.Request) bool {
		return context(r).DefaultDomain
	}
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
