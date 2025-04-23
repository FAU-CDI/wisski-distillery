// Package static implements serving of fully static resources
//
//spellchecker:words assets
package assets

//spellchecker:words context embed http github wisski distillery internal component
import (
	"context"
	"embed"
	"fmt"
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
		Prefix: Public,

		CSRF: false,
	}
}

//go:embed dist
var staticFS embed.FS

func (static *Static) HandleRoute(ctx context.Context, route string) (http.Handler, error) {
	// take the filesystem
	fs, err := fs.Sub(staticFS, "dist")
	if err != nil {
		return nil, fmt.Errorf("failed to get 'dist' directory: %w", err)
	}

	// and serve it
	return http.StripPrefix(route, http.FileServer(http.FS(fs))), nil
}
