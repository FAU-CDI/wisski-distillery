// Package static implements serving of fully static resources
package static

import (
	"context"
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/tkw1536/goprogram/stream"
)

type Static struct {
	component.ComponentBase
}

func (*Static) Name() string { return "static" }

func (*Static) Routes() []string { return []string{"/static/"} }

func (static *Static) Handler(route string, context context.Context, io stream.IOStream) (http.Handler, error) {
	fs, err := fs.Sub(distStaticFS, "dist")
	if err != nil {
		return nil, err
	}

	// censor *.html in the filesystem
	fs = fsx.Censor(fs, func(path string) bool {
		suffix := "html"
		return len(path) >= len(suffix) && strings.EqualFold(path[len(path)-len(suffix):], suffix)
	})

	// and serve it
	return http.StripPrefix(route, http.FileServer(http.FS(fs))), nil
}

//go:embed dist
var distStaticFS embed.FS
