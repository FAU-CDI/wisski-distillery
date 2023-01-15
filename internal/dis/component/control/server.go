package control

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/cancel"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/FAU-CDI/wisski-distillery/pkg/mux"
	"github.com/gorilla/csrf"
	"github.com/rs/zerolog"
)

// Server returns an http.Mux that implements the main server instance.
// The server may spawn background tasks, but these should be terminated once context closes.
//
// Logging messages are directed to progress
func (control *Control) Server(ctx context.Context, progress io.Writer) (public http.Handler, internal http.Handler, err error) {
	logger := zerolog.Ctx(ctx)

	var publicM, internalM mux.Mux[component.RouteContext]
	publicM.Context = func(r *http.Request) component.RouteContext {
		slug, ok := control.Still.Config.SlugFromHost(r.Host)
		return component.RouteContext{
			DefaultDomain: slug == "" && ok,
		}
	}
	publicM.Panic = func(panic any, w http.ResponseWriter, r *http.Request) {
		// log the panic
		logger.Error().
			Str("panic", fmt.Sprint(panic)).
			Str("path", r.URL.Path).
			Msg("panic serving handler")

		// and send an internal server error
		httpx.TextInterceptor.Fallback.ServeHTTP(w, r)
	}

	// setup the internal server identically
	internalM.Panic = publicM.Panic
	internalM.Context = publicM.Context

	// create a csrf protector
	csrfProtector := control.CSRF()

	// iterate over all the handler
	for _, s := range control.Dependencies.Routeables {
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
		predicate := routes.Predicate(publicM.ContextOf)

		// and add all the prefixes
		for _, prefix := range append([]string{routes.Prefix}, routes.Aliases...) {
			if routes.Internal {
				internalM.Add(prefix, predicate, routes.Exact, handler)
			} else {
				publicM.Add(prefix, predicate, routes.Exact, handler)
			}
		}
	}

	// apply the given context function
	public = httpx.WithContextWrapper(&publicM, func(rcontext context.Context) context.Context { return cancel.ValuesOf(rcontext, ctx) })
	internal = httpx.WithContextWrapper(&internalM, func(rcontext context.Context) context.Context { return cancel.ValuesOf(rcontext, ctx) })
	err = nil
	return
}

// CSRF returns a CSRF handler for the given function
func (control *Control) CSRF() func(http.Handler) http.Handler {
	var opts []csrf.Option
	if !control.Config.HTTPSEnabled() {
		opts = append(opts, csrf.Secure(false))
	}
	opts = append(opts, csrf.SameSite(csrf.SameSiteStrictMode))
	opts = append(opts, csrf.CookieName(CSRFCookie))
	opts = append(opts, csrf.FieldName(CSRFCookieField))
	return csrf.Protect(control.Config.CSRFSecret(), opts...)
}
