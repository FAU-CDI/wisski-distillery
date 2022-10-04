package component

import (
	"net/http"

	"github.com/tkw1536/goprogram/stream"
)

// Servable implements a component with a Serve method
type Servable interface {
	Component

	// Routes returns the routes served by this servable
	Routes() []string

	// Handler returns the handler for the requested route
	Handler(route string, io stream.IOStream) (http.Handler, error)
}
