package server

import (
	"context"
	"io"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/handling"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/tkw1536/pkglib/contextx"
	"github.com/tkw1536/pkglib/httpx/mux"
	"github.com/tkw1536/pkglib/httpx/wrap"
	"github.com/tkw1536/pkglib/recovery"

	"github.com/gorilla/csrf"
	"github.com/rs/zerolog"
)

// Server represents the running control server.
type Server struct {
	component.Base
	dependencies struct {
		Routeables []component.Routeable
		Cronables  []component.Cronable

		Templating *templating.Templating
		Handleing  *handling.Handling
	}
}

var (
	_ component.Installable = (*Server)(nil)
)

// Server returns an http.Mux that implements the main server instance.
// The server may spawn background tasks, but these should be terminated once context closes.
//
// Logging messages are directed to progress
func (server *Server) Server(ctx context.Context, progress io.Writer) (public http.Handler, internal http.Handler, err error) {
	interceptor := server.dependencies.Handleing.TextInterceptor()

	// wrapHandler wraps individual handlers for errors
	wrapHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// handle any panic()s that occur
			defer func() {
				// intercept any panic() that wasn't caught
				if err := recovery.Recover(recover()); err != nil {
					interceptor.Intercept(w, r, err)
				}
			}()

			// determine if we are on a slug from a host
			slug, ok := server.Config.HTTP.NormSlugFromHost(r.Host)

			rctx := component.WithRouteContext(r.Context(), component.RouteContext{
				DefaultDomain: slug == "" && ok,
			})
			ctx := contextx.WithValuesOf(rctx, ctx)

			// serve with the next context
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	var publicM, internalM mux.Mux

	// create a csrf protector
	csrfProtector := server.csrf()

	// iterate over all the handler
	for _, s := range server.dependencies.Routeables {
		routes := s.Routes()
		zerolog.Ctx(ctx).Info().
			Str("Name", s.Name()).
			Str("Prefix", routes.Prefix).
			Strs("Aliases", routes.Aliases).
			Bool("Exact", routes.Exact).
			Bool("CSRF", routes.CSRF).
			Bool("Decorator", routes.Decorator != nil).
			Bool("Internal", routes.Internal).
			Bool("MatchAllDomains", routes.MatchAllDomains).
			Msg("mounting route")

		// call the handler for the route
		handler, err := s.HandleRoute(ctx, routes.Prefix)
		if err != nil {
			zerolog.Ctx(ctx).Err(err).
				Str("Component", s.Name()).
				Str("Prefix", routes.Prefix).
				Msg("error mounting route")
			continue
		}

		// decorate the handler
		handler = routes.Decorate(handler, csrfProtector)

		// determine the predicate
		predicate := routes.Predicate(func(r *http.Request) component.RouteContext { return component.RouteContextOf(r.Context()) })

		// and add all the prefixes
		for _, prefix := range append([]string{routes.Prefix}, routes.Aliases...) {
			if routes.Internal {
				internalM.Add(prefix, predicate, routes.Exact, handler)
			} else {
				publicM.Add(prefix, predicate, routes.Exact, handler)
			}
		}
	}

	// wrap the handlers
	public = wrapHandler(&publicM)
	internal = wrapHandler(&internalM)

	// Add Content-Security-Policy
	public = WithCSP(public, models.ContentSecurityPolicyDistilery)
	internal = WithCSP(internal, models.ContentSecurityPolicyNothing)

	public = wrap.Time(public)

	err = nil
	return
}

// CSRF returns a CSRF handler for the given function
func (server *Server) csrf() func(http.Handler) http.Handler {
	var opts []csrf.Option
	opts = append(opts, csrf.Secure(server.Config.HTTP.HTTPSEnabled()))
	opts = append(opts, csrf.SameSite(csrf.SameSiteStrictMode))
	opts = append(opts, csrf.CookieName(CSRFCookie))
	opts = append(opts, csrf.FieldName(CSRFCookieField))
	return csrf.Protect(server.Config.CSRFSecret(), opts...)
}

// WithCSP adds a Content-Security-Policy header to every response
func WithCSP(handler http.Handler, policy string) http.Handler {
	if policy == "" {
		return handler
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", policy)
		handler.ServeHTTP(w, r)
	})
}
