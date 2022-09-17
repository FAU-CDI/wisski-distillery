package control

import (
	"net/http"

	"github.com/tkw1536/goprogram/stream"
)

// Server returns an http.Mux that implements the main server instance
func (control Control) Server(io stream.IOStream) (http.Handler, error) {
	// self server
	self, err := control.self(io)
	if err != nil {
		return nil, err
	}

	resolver, err := control.resolver(io)
	if err != nil {
		return nil, err
	}

	info, err := control.info(io)
	if err != nil {
		return nil, err
	}

	// resolver

	mux := http.NewServeMux()
	mux.Handle("/", self)

	mux.Handle("/go/", resolver)
	mux.Handle("/wisski/get/", resolver)

	// TODO: Fix me!
	mux.Handle("/dis/", info)

	return mux, nil
}
