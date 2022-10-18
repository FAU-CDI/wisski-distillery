// Package static implements serving of fully static resources
package static

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/tkw1536/goprogram/stream"
)

type Static struct {
	component.Base
}

func (*Static) Routes() []string { return []string{"/static/"} }

//go:embed dist
var staticFS embed.FS

func (static *Static) Handler(route string, context context.Context, io stream.IOStream) (http.Handler, error) {
	// take the filesystem
	fs, err := fs.Sub(staticFS, "dist")
	if err != nil {
		return nil, err
	}

	// and serve it
	return http.StripPrefix(route, http.FileServer(http.FS(fs))), nil
}
