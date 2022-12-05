package control

import (
	"context"
	"io"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/pkg/cancel"
	"github.com/rs/zerolog"
)

// Server returns an http.Mux that implements the main server instance.
// The server may spawn background tasks, but these should be terminated once context closes.
//
// Logging messages are directed to progress
func (control *Control) Server(ctx context.Context, progress io.Writer) (http.Handler, error) {
	// create a new mux
	mux := http.NewServeMux()

	// add all the servable routes!
	for _, s := range control.Dependencies.Routeables {
		for _, route := range s.Routes() {
			zerolog.Ctx(ctx).Info().Str("component", s.Name()).Str("route", route).Msg("mounting route")
			handler, err := s.HandleRoute(ctx, route)
			if err != nil {
				return nil, err
			}
			mux.Handle(route, handler)
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(cancel.ValuesOf(r.Context(), ctx))
		mux.ServeHTTP(w, r)
	}), nil
}
