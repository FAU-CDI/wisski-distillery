//spellchecker:words server
package server

//spellchecker:words context http github wisski distillery internal component server handling templating models wdlog pkglib contextx httpx wrap recovery gorilla csrf
import (
	"context"
	"io"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/handling"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/tkw1536/pkglib/contextx"
	"github.com/tkw1536/pkglib/httpx/mux"
	"github.com/tkw1536/pkglib/httpx/wrap"
	"github.com/tkw1536/pkglib/recovery"

	"github.com/gorilla/csrf"
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
			slug, ok := component.GetStill(server).Config.HTTP.NormSlugFromHost(r.Host)

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
		wdlog.Of(ctx).Info(
			"mounting route",

			"Name", s.Name(),
			"Prefix", routes.Prefix,
			"Aliases", routes.Aliases,
			"Exact", routes.Exact,
			"CSRF", routes.CSRF,
			"Decorator", routes.Decorator != nil,
			"Internal", routes.Internal,
			"MatchAllDomains", routes.MatchAllDomains,
		)

		// call the handler for the route
		handler, err := s.HandleRoute(ctx, routes.Prefix)
		if err != nil {
			wdlog.Of(ctx).Error(
				"error mounting route",
				"error", err,

				"Component", s.Name(),
				"Prefix", routes.Prefix,
			)
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
	public = WithCSP(public, models.ContentSecurityPolicyPanel)
	internal = WithCSP(internal, models.ContentSecurityPolicyNothing)

	public = wrap.Time(public)

	err = nil
	return
}

// CSRF returns a CSRF handler for the given function
func (server *Server) csrf() func(http.Handler) http.Handler {
	config := component.GetStill(server).Config

	var opts []csrf.Option
	opts = append(opts, csrf.Secure(config.HTTP.HTTPSEnabled()))
	opts = append(opts, csrf.SameSite(csrf.SameSiteStrictMode))
	opts = append(opts, csrf.Path("/"))
	opts = append(opts, csrf.CookieName(CSRFCookie))
	opts = append(opts, csrf.FieldName(CSRFCookieField))
	return csrf.Protect(config.CSRFSecret(), opts...)
}

// WithCSP adds a Content-Security-Policy header to every response
func WithCSP(handler http.Handler, policy string) http.Handler {
	if policy == "" {
		return handler
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetCSP(w, policy)
		handler.ServeHTTP(w, r)
	})
}

const cspHeader = "Content-Security-Policy"

// SetCSP sets the Content-Security-Policy for the given response
// Any previously set header is discarded
func SetCSP(w http.ResponseWriter, policy string) {
	header := w.Header()

	header.Del(cspHeader)
	header.Set(cspHeader, policy)
}
