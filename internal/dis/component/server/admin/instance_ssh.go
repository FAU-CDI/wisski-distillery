package admin

import (
	"context"
	_ "embed"
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/pkglib/httpx"

	"github.com/julienschmidt/httprouter"

	gossh "golang.org/x/crypto/ssh"
)

//go:embed "html/instance_ssh.html"
var instanceSSHHTML []byte
var instanceSSHTemplate = templating.Parse[instanceSSHContext](
	"instance_ssh.html", instanceSSHHTML, nil,

	templating.Assets(assets.AssetsAdmin),
)

type instanceSSHContext struct {
	templating.RuntimeFlags

	Instance *wisski.WissKI
	SSHKeys  []string

	Hostname    string
	PanelDomain string
	Port        uint16
}

func (admin *Admin) instanceSSH(ctx context.Context) http.Handler {
	tpl := instanceSSHTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuInstance,
			menuSSH,
		),
	)

	return tpl.HTMLHandlerWithFlags(func(r *http.Request) (ctx instanceSSHContext, funcs []templating.FlagFunc, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		// setup the context with just the instance
		ctx.Instance, err = admin.dependencies.Instances.WissKI(r.Context(), slug)
		if err != nil {
			return ctx, nil, httpx.ErrNotFound
		}

		ctx.Hostname = ctx.Instance.Domain()
		ctx.PanelDomain = admin.Config.HTTP.PanelDomain()
		ctx.Port = admin.Config.Listen.SSHPort

		keys, err := ctx.Instance.SSH().Keys(r.Context())
		if err != nil {
			return ctx, nil, httpx.ErrInternalServerError
		}

		ctx.SSHKeys = make([]string, len(keys))
		for i, key := range keys {
			ctx.SSHKeys[i] = string(gossh.MarshalAuthorizedKey(key))
		}

		return ctx, []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + ctx.Instance.Slug)}),
			templating.ReplaceCrumb(menuSSH, component.MenuItem{Title: "SSH", Path: template.URL("/admin/instance/" + ctx.Instance.Slug + "/ssh")}),
			templating.Title(ctx.Instance.Slug + " - SSH"),
			admin.instanceTabs(slug, "ssh"),
		}, nil
	})
}
