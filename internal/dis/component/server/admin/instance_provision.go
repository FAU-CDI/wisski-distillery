package admin

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/models"

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

	systemParams
}

type systemParams struct {
	PHPVersions             []string
	ContentSecurityPolicies []string
	DefaultPHPVersion       string
}

func newSystemParams() (sp systemParams) {
	sp.PHPVersions = models.KnownPHPVersions()
	sp.ContentSecurityPolicies = models.ContentSecurityPolicyExamples()
	sp.DefaultPHPVersion = models.DefaultPHPVersion
	return sp
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
		ipc.systemParams = newSystemParams()
		return ipc, nil
	})
}
