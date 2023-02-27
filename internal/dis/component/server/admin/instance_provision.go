package admin

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"

	_ "embed"
)

//go:embed "html/instance_provision.html"
var instanceProvisionHTML []byte
var instanceProvisionTemplate = templating.Parse[instanceProvisionContext](
	"instance_provision.html", instanceProvisionHTML, nil,

	templating.Title("Provision New Instance"),
	templating.Assets(assets.AssetsAdminProvision),
)

type instanceProvisionContext struct {
	templating.RuntimeFlags

	// nothing for the moment
}

func (admin *Admin) instanceProvision(ctx context.Context) http.Handler {
	tpl := instanceProvisionTemplate.Prepare(
		admin.Dependencies.Templating,

		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuProvision,
		),
	)

	return tpl.HTMLHandler(func(r *http.Request) (ipc instanceProvisionContext, err error) {
		return ipc, nil
	})
}
