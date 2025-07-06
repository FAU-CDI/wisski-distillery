//spellchecker:words socket
package socket

//spellchecker:words context http strings github process over websocket proto wisski distillery internal component auth scopes exporter instances purger provision server admin socket actions models pkglib lazy
import (
	"context"
	"net/http"
	"strings"

	"github.com/FAU-CDI/process_over_websocket"
	"github.com/FAU-CDI/process_over_websocket/proto"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances/purger"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/provision"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/admin/socket/actions"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"go.tkw01536.de/pkglib/lazy"
)

type Sockets struct {
	component.Base

	handler lazy.Lazy[proto.Handler]

	dependencies struct {
		Actions  []actions.WebsocketAction
		IActions []actions.WebsocketInstanceAction

		Provision *provision.Provision
		Instances *instances.Instances
		Exporter  *exporter.Exporter
		Purger    *purger.Purger
		Auth      *auth.Auth
	}
}

var (
	_ component.Routeable = (*Sockets)(nil)
)

func (socket *Sockets) Routes() component.Routes {
	return component.Routes{
		Prefix:    "/api/v1/pow",
		Decorator: socket.dependencies.Auth.Require(true, scopes.ScopeUserValid, nil),
	}
}

func (sockets *Sockets) HandleRoute(ctx context.Context, path string) (http.Handler, error) {
	pow := process_over_websocket.Server{
		Handler: sockets.handler.Get(func() proto.Handler { return sockets.Actions(ctx) }),
		Options: process_over_websocket.Options{
			BasePath: "/api/v1/pow/",
		},
	}
	pow.Options.RESTOptions.OpenAPIServerDescription = "Distillery POW Server"

	// ensure that the server is closed once we are
	go func() {
		<-ctx.Done()
		pow.Close()
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if we're in the docs, unsafely set the unsafe csp
		if strings.HasPrefix(r.URL.Path, "/api/v1/pow/docs/") {
			server.SetCSP(w, models.ContentSecurityPolicyPanelUnsafeScripts)
		}

		pow.ServeHTTP(w, r)
	}), nil
}
