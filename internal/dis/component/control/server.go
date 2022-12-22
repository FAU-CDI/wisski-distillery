package control

import (
	"context"
	"io"
	"net/http"

	"github.com/rs/zerolog"
)

// Server returns an http.Mux that implements the main server instance.
// The server may spawn background tasks, but these should be terminated once context closes.
//
// Logging messages are directed to progress
func (control *Control) Server(ctx context.Context, progress io.Writer) (*http.ServeMux, error) {
	// create a new mux
	mux := http.NewServeMux()

	// add all the servable routes!
	for _, s := range control.Dependencies.Servables {
		for _, route := range s.Routes() {
			zerolog.Ctx(ctx).Info().Str("component", s.Name()).Str("route", route).Msg("mounting route")
			handler, err := s.Handler(ctx, route)
			if err != nil {
				return nil, err
			}
			mux.Handle(route, handler)
		}
	}
	return mux, nil
}
