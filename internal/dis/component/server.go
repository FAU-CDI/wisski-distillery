package component

import (
	"context"
	"io"
	"net/http"
)

// Servable is a component that is servable
type Servable interface {
	Component

	// Routes returns the routes served by this servable
	Routes() []string

	// Handler returns the handler for the requested route
	Handler(ctx context.Context, route string, progress io.Writer) (http.Handler, error)
}
