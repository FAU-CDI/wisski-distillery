package control

import (
	"context"
	"net/http"

	"github.com/tkw1536/goprogram/stream"
)

// Server returns an http.Mux that implements the main server instance.
// The server may spawn background tasks, but these should be terminated once context closes.
//
// Logging messages are directed to io.
func (control *Control) Server(ctx context.Context, io stream.IOStream) (*http.ServeMux, error) {
	// create a new mux
	mux := http.NewServeMux()

	// add all the servable routes!
	for _, s := range control.Servables {
		for _, route := range s.Routes() {
			io.Printf("mounting %s\n", route)
			handler, err := s.Handler(ctx, route, io)
			if err != nil {
				return nil, err
			}
			mux.Handle(route, handler)
		}
	}
	return mux, nil
}
