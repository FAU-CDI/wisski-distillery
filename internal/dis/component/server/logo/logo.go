//spellchecker:words logo
package logo

//spellchecker:words context http github wisski distillery internal component server models pkglib httpx embed
import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/tkw1536/pkglib/httpx"

	_ "embed"
)

type Logo struct {
	component.Base
}

var (
	_ component.Routeable = (*Logo)(nil)
)

func (*Logo) Routes() component.Routes {
	return component.Routes{
		Prefix:  "/logo/",
		Aliases: []string{"/favicon.ico", "/logo.svg"},
		Exact:   true,
	}
}

var (
	//go:embed favicon.ico
	faviconICO []byte

	//go:embed logo.svg
	logoSVG []byte
)

var faviconRoute = httpx.Response{
	ContentType: "image/x-icon",
	Body:        faviconICO,
}

var logoSVGRoute = httpx.Response{
	ContentType: "image/svg+xml",
	Body:        logoSVG,
}.Minify()

func (*Logo) HandleRoute(ctx context.Context, path string) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/favicon.ico":
			faviconRoute.ServeHTTP(w, r)
		case "/logo.svg":
			server.SetCSP(w, models.ContentSecurityPolicyPanelUnsafeStyles)
			logoSVGRoute.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	}), nil
}
