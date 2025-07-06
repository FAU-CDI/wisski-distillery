//spellchecker:words component
package component

//spellchecker:words context http github pkglib httpx
import (
	"context"
	"net/http"

	"go.tkw01536.de/pkglib/httpx/mux"
)

// Routeable is a component that is servable.
type Routeable interface {
	Component

	// Routes returns information about the routes to be handled by this Routeable
	Routes() Routes

	// HandleRoute returns the handler for the requested path.
	// Context is cancelled once the handler should be closed.
	HandleRoute(ctx context.Context, path string) (http.Handler, error)
}

// Routes represents information about a single Routeable.
type Routes struct {
	// Prefix is the prefix this pattern handles
	Prefix string

	// MatchAllDomains indicates that all domains, even the non-default domain, should be matched
	MatchAllDomains bool

	// Internal indicates that this route should only answer on the internal server.
	// Internal implies MatchAllDomains.
	Internal bool

	// MenuTitle, MenuSticky and MenuPriority return the priority, sticky and title of this menu item
	// see MenuItem for details
	MenuTitle    string
	MenuSticky   bool
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

type routeContextTyp int

const routeContextKey routeContextTyp = 0

// RouteContext represents the context passed to a given route.
type RouteContext struct {
	DefaultDomain bool
}

// WithRouteContext adds the given RouteContext to the context.
func WithRouteContext(parent context.Context, value RouteContext) context.Context {
	return context.WithValue(parent, routeContextKey, value)
}

// RouteContextOf returns the route context of the given context.
func RouteContextOf(context context.Context) RouteContext {
	ctx, ok := context.Value(routeContextKey).(RouteContext)
	if !ok {
		return RouteContext{}
	}
	return ctx
}

// Predicate returns the predicate corresponding to the given route.
func (routes Routes) Predicate(context func(*http.Request) RouteContext) mux.Predicate {
	if routes.MatchAllDomains || routes.Internal {
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
