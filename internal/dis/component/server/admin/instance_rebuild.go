//spellchecker:words admin
package admin

//spellchecker:words context html template http github wisski distillery internal component instances server assets templating models julienschmidt httprouter pkglib httpx embed
import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/julienschmidt/httprouter"
	"github.com/tkw1536/pkglib/httpx"

	_ "embed"
)

//go:embed "html/instance_system.html"
var instanceSystemHTML []byte
var instanceSystemTemplate = templating.Parse[instanceSystemContext](
	"instance_system.html", instanceSystemHTML, nil,
)

// instanceSystemContext is the context for instance_system.html.
type instanceSystemContext struct {
	templating.RuntimeFlags

	// parameters for completion
	PHPVersions             []string
	ContentSecurityPolicies []string
	DefaultPHPVersion       string

	// Are we in rebuild mode?
	Rebuild bool
	Slug    string
	System  models.System

	// list of known profiles and their descriptions
	DefaultProfile string
	Profiles       map[string]string
}

// prepare prares the given instanceSystemContent.
func (isc *instanceSystemContext) prepare(rebuild bool) {
	isc.Rebuild = rebuild
	isc.PHPVersions = models.KnownPHPVersions()
	isc.ContentSecurityPolicies = models.ContentSecurityPolicyExamples()
	isc.DefaultPHPVersion = models.DefaultPHPVersion
}

func (admin *Admin) instanceRebuild(context.Context) http.Handler {
	tpl := instanceSystemTemplate.Prepare(
		admin.dependencies.Templating,

		templating.Title("Rebuild Instance"),
		templating.Assets(assets.AssetsAdminRebuild),

		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuInstance,
			menuRebuild,
		),
	)

	return tpl.HTMLHandlerWithFlags(admin.dependencies.Handling, func(r *http.Request) (isc instanceSystemContext, funcs []templating.FlagFunc, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		var instance *wisski.WissKI
		instance, err = admin.dependencies.Instances.WissKI(r.Context(), slug)
		if errors.Is(err, instances.ErrWissKINotFound) {
			return isc, nil, httpx.ErrNotFound
		}
		if err != nil {
			return isc, nil, fmt.Errorf("failed to get WissKI: %w", err)
		}

		isc.Slug = instance.Slug
		isc.System = instance.System

		// replace the menu item
		escapedSlug := url.PathEscape(instance.Slug)
		funcs = []templating.FlagFunc{
			templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + escapedSlug)}),            // #nosec G203 -- escaped and safe
			templating.ReplaceCrumb(menuRebuild, component.MenuItem{Title: "Rebuild", Path: template.URL("/admin/instance/" + escapedSlug + "/rebuild")}), // #nosec G203 -- escaped and safe
			templating.Title(instance.Slug + " - Rebuild"),
			admin.instanceTabs(escapedSlug, "rebuild"),
		}

		isc.prepare(true)
		return
	})
}
