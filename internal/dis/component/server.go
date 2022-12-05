package component

import (
	"context"
	"net/http"
)

// Routeable is a component that is servable
type Routeable interface {
	Component

	// Routes returns the routes served by this servable
	Routes() []string

	// HandleRoute returns the handler for the requested route
	HandleRoute(ctx context.Context, route string) (http.Handler, error)
}
