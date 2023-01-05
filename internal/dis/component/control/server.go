package control

import (
	"context"
	"io"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/pkg/cancel"
	"github.com/gorilla/csrf"
	"github.com/rs/zerolog"
)

// Server returns an http.Mux that implements the main server instance.
// The server may spawn background tasks, but these should be terminated once context closes.
//
// Logging messages are directed to progress
func (control *Control) Server(ctx context.Context, progress io.Writer) (http.Handler, error) {
	// create a new mux
	mux := http.NewServeMux()

	// create a csrf protector
	csrfProtector := control.CSRF()

	// iterate over all the handler
	for _, s := range control.Dependencies.Routeables {
		routes := s.Routes()
		zerolog.Ctx(ctx).Info().Str("component", s.Name()).Strs("paths", routes.Paths).Bool("csrf", routes.CSRF).Bool("decorator", routes.Decorator != nil).Msg("mounting route")

		for _, path := range routes.Paths {
			handler, err := s.HandleRoute(ctx, path)
			if err != nil {
				return nil, err
			}
			mux.Handle(path, routes.Decorate(handler, csrfProtector))
		}
	}

	// apply the given context function
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(cancel.ValuesOf(r.Context(), ctx))
		mux.ServeHTTP(w, r)
	}), nil
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
