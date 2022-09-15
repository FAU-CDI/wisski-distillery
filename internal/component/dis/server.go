package dis

import (
	"net/http"

	"github.com/tkw1536/goprogram/stream"
)

// Server returns an http.Mux that implements the main server instance
func (dis Dis) Server(io stream.IOStream) (http.Handler, error) {
	// self server
	self, err := dis.self(io)
	if err != nil {
		return nil, err
	}

	resolver, err := dis.resolver(io)
	if err != nil {
		return nil, err
	}

	info, err := dis.info(io)
	if err != nil {
		return nil, err
	}

	// resolver

	mux := http.NewServeMux()
	mux.Handle("/", self)

	mux.Handle("/go/", resolver)
	mux.Handle("/wisski/navigate", resolver)

	// TODO: Fix me!
	mux.Handle("/dis/", info)

	return mux, nil
}
