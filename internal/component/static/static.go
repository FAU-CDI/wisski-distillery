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
	component.ComponentBase
}

func (*Static) Name() string { return "static" }

func (*Static) Routes() []string { return []string{"/static/"} }

func (static *Static) Handler(route string, context context.Context, io stream.IOStream) (http.Handler, error) {
	fs, err := fs.Sub(htmlStaticFS, "out")
	if err != nil {
		return nil, err
	}

	return http.StripPrefix(route, http.FileServer(http.FS(fs))), nil
}

//go:embed out
var htmlStaticFS embed.FS
