package control

import (
	"net/http"

	"github.com/tkw1536/goprogram/stream"
)

// Server returns an http.Mux that implements the main server instance
// Logging messages are directed to io.
func (control *Control) Server(io stream.IOStream) (*http.ServeMux, error) {
	// create a new mux
	mux := http.NewServeMux()

	// add all the servable routes!
	for _, s := range control.Servables {
		for _, route := range s.Routes() {
			io.Printf("mounting %s\n", route)
			handler, err := s.Handler(route, io)
			if err != nil {
				return nil, err
			}
			mux.Handle(route, handler)
		}
	}
	return mux, nil
}
