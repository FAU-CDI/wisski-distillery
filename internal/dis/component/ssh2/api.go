package ssh2

//spellchecker:words context http github wisski distillery internal component pkglib httpx golang crypto gossh
import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/tkw1536/pkglib/httpx"
	gossh "golang.org/x/crypto/ssh"
)

func (ssh2 *SSH2) Routes() component.Routes {
	return component.Routes{
		Prefix:   "/authorized_keys/",
		Exact:    true,
		Internal: true,
	}
}
func (ssh2 *SSH2) HandleRoute(ctx context.Context, path string) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fetch the global keys
		gkeys, err := ssh2.dependencies.Keys.Admin(r.Context())
		if err != nil {
			httpx.TextInterceptor.Intercept(w, r, err)
			return
		}

		// find the host
		slug, ok := component.GetStill(ssh2).Config.HTTP.SlugFromHost(r.Host)
		if slug == "" || !ok {
			httpx.TextInterceptor.Intercept(w, r, httpx.ErrNotFound)
			return
		}

		// fetch the instance
		instance, err := ssh2.dependencies.Instances.WissKI(r.Context(), slug)
		if err != nil {
			httpx.TextInterceptor.Intercept(w, r, httpx.ErrNotFound)
			return
		}

		// fetch the instance keys
		keys, err := instance.SSH().Keys(r.Context())
		if err != nil {
			httpx.TextInterceptor.Intercept(w, r, err)
			return
		}

		// marshal out everything!
		for _, key := range gkeys {
			w.Write(gossh.MarshalAuthorizedKey(key))
		}
		for _, key := range keys {
			w.Write(gossh.MarshalAuthorizedKey(key))
		}
	}), nil
}
