// Package static implements serving of fully static resources
package static

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

type Static struct {
	component.Base
}

var (
	_ component.Routeable = (*Static)(nil)
)

func (*Static) Routes() component.Routes {
	return component.Routes{
		Paths: []string{"/static/"},
		CSRF:  false,
	}
}

//go:embed dist
var staticFS embed.FS

func (static *Static) HandleRoute(ctx context.Context, route string) (http.Handler, error) {
	// take the filesystem
	fs, err := fs.Sub(staticFS, "dist")
	if err != nil {
		return nil, err
	}

	// and serve it
	return http.StripPrefix(route, http.FileServer(http.FS(fs))), nil
}
