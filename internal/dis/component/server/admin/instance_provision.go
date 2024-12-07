//spellchecker:words admin
package admin

//spellchecker:words context http github wisski distillery internal component server assets templating ingredient barrel manager pkglib collection embed
import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/manager"
	"github.com/tkw1536/pkglib/collection"

	_ "embed"
)

func (admin *Admin) instanceProvision(context.Context) http.Handler {
	tpl := instanceSystemTemplate.Prepare(
		admin.dependencies.Templating,

		templating.Title("Provision New Instance"),
		templating.Assets(assets.AssetsAdminProvision),

		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuProvision,
		),
	)

	return tpl.HTMLHandler(admin.dependencies.Handling, func(r *http.Request) (ipc instanceSystemContext, err error) {
		ipc.prepare(false)
		ipc.DefaultProfile = manager.DefaultProfile()
		ipc.Profiles = collection.MapValues(manager.Profiles(), func(_ string, profile manager.Profile) string { return profile.Description })
		return ipc, nil
	})
}
